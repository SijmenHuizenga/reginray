package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/plugins/logdriver"
	"github.com/docker/docker/daemon/logger"
	protoio "github.com/gogo/protobuf/io"
	"github.com/pkg/errors"
	"github.com/tonistiigi/fifo"
	"time"
	"strings"
)

type Driver struct {
	mutex         sync.Mutex
	logFiles      map[string]*Logger
	logContainers map[string]*Logger
}

type Logger struct {
	stream io.ReadCloser
	info   logger.Info
}

func newDriver() *Driver {
	return &Driver{
		logFiles:      make(map[string]*Logger),
		logContainers: make(map[string]*Logger),
	}
}

func (driver *Driver) StartLogging(file string, logContext logger.Info) error {
	driver.mutex.Lock()
	if _, exists := driver.logFiles[file]; exists {
		driver.mutex.Unlock()
		return fmt.Errorf("logger for %q already exists", file)
	}
	driver.mutex.Unlock()

	if logContext.LogPath == "" {
		logContext.LogPath = filepath.Join("/var/log/docker", logContext.ContainerID)
	}
	if err := os.MkdirAll(filepath.Dir(logContext.LogPath), 0755); err != nil {
		return errors.Wrap(err, "error setting up logger dir")
	}

	logrus.WithField("id", logContext.ContainerID).WithField("file", file).
		WithField("logpath", logContext.LogPath).Debugf("Start logging")
	fifoLogFile, err := fifo.OpenFifo(context.Background(), file, syscall.O_RDONLY, 0700)
	if err != nil {
		return errors.Wrapf(err, "error opening logger fifo: %q", file)
	}

	driver.mutex.Lock()
	lf := &Logger{fifoLogFile, logContext}
	driver.logFiles[file] = lf
	driver.logContainers[logContext.ContainerID] = lf
	driver.mutex.Unlock()

	go consumeLogs(lf)
	return nil
}

func (driver *Driver) StopLogging(file string) error {
	logrus.WithField("file", file).Debugf("Stop logging")
	driver.mutex.Lock()
	lf, ok := driver.logFiles[file]
	if ok {
		lf.stream.Close()
		delete(driver.logContainers, driver.logFiles[file].info.ContainerID)
		delete(driver.logFiles, file)
	}
	driver.mutex.Unlock()
	return nil
}

func consumeLogs(lf *Logger) {
	// create a protobuf reader for the log stream
	dec := protoio.NewUint32DelimitedReader(lf.stream, binary.BigEndian, 1e6)
	defer dec.Close()
	defer lf.stream.Close()
	// a temp buffer for each log entry
	var buf logdriver.LogEntry

	for {
		// reads a message from the log stream and put it in a buffer
		readFailCounter := 0
		for readErr := dec.ReadMsg(&buf); readErr != nil; {
			// exit the loop if reader reaches EOF or the fifo is closed by the writer
			if readErr == io.EOF || readErr == os.ErrClosed || strings.Contains(readErr.Error(), "file already closed") {
				logrus.WithField("id", lf.info.ContainerID).WithError(readErr).Info("shutting down loggers")
				return
			}

			// exit the loop if retry number reaches the specified number
			if readFailCounter > 10 {
				logrus.WithField("id", lf.info.ContainerID).WithField("readFailCounter", readFailCounter).
					WithError(readErr).Error("Stop retrying. Shutting down loggers")
				return
			}

			// if there is any other error, retry for robustness.
			readFailCounter++
			logrus.WithField("id", lf.info.ContainerID).WithField("readFailCounter", readFailCounter).
				WithError(readErr).Error("Encountered error and retrying")
			time.Sleep(1 * time.Second)

			dec = protoio.NewUint32DelimitedReader(lf.stream, binary.BigEndian, 1e6)
		}

		sendLogs(lf, buf)

		buf.Reset()
	}
}

func sendLogs(lf *Logger, buf logdriver.LogEntry) {
	line := string(buf.Line[:])
	sendFailCount := 0
	for sendErr := sendLogToAspicio(lf.info, line, buf.TimeNano); sendErr != nil; {
		sendFailCount++

		if sendFailCount >= 10 {
			//todo: store logs that could not be sent on persistant disk so they can be sent later.
			logrus.WithField("id", lf.info.ContainerID).WithField("sendFailCount", sendFailCount).
				WithError(sendErr).WithField("line", line).Error("Failed Stop retrying sending to aspicio backend. Skipping this line")
			return
		}
		time.Sleep(1 * time.Second)
	}
}
