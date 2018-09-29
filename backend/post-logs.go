package main

import (
	e "./http"
	"net/http"
	"strconv"
	"./model"
	"encoding/json"
)

type InputLogEntry struct {
	Message   string
	Timestamp uint64
	Container InputLogEntryContainer
	Image     InputLogEntryImage
}

type InputLogEntryContainer struct {
	Name   string
	Id     string
	Labels map[string]string
}

type InputLogEntryImage struct {
	Id   string
	Name string
}

func (a *Aspicio) PostLogs(writer http.ResponseWriter, request *http.Request) {
	var inputlogentry InputLogEntry
	err := json.NewDecoder(request.Body).Decode(&inputlogentry)
	if err != nil {
		e.FailOnError(err, writer)
		return
	}

	logentry := model.LogEntry{
		Timestamp:        strconv.FormatUint(inputlogentry.Timestamp, 10),
		TimestampSeconds: uint32(inputlogentry.Timestamp / uint64(1000000000)),
		Fields:           map[string]string{},
		Message:          inputlogentry.Message,
		Image:            inputlogentry.Image,
		Container:        inputlogentry.Container,
	}

	a.decoder.AddFields(&logentry)

	err = a.logsdb.Collection.Insert(logentry)
}
