package main

import (
	"github.com/docker/docker/daemon/logger"
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
)

type LogEntry struct {
	Message   string
	Timestamp int64
	Container LogEntryContainer
	Image     LogEntryImage
}

type LogEntryContainer struct {
	Name   string
	Id     string
	Labels map[string]string
}

type LogEntryImage struct {
	Id   string
	Name string
}

func sendLogToAspicio(info logger.Info, line string, timestamp int64) error {
	url := backendUrl + "/logs"
	entry := LogEntry{
		Message:   line,
		Timestamp: timestamp,
		Container: LogEntryContainer{
			Name:   info.ContainerName,
			Id:     info.ContainerID,
			Labels: info.ContainerLabels,
		},
		Image: LogEntryImage{
			Id:   info.ContainerImageID,
			Name: info.ContainerImageName,
		},
	}

	jsonValue, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		bodyStr := buf.String()
		resp.Body.Close()

		return errors.New("Response code from aspicio is not 200. Body: " + bodyStr)
	}

	return nil
}
