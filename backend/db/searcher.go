package db

import (
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"net/http"
	"github.com/pkg/errors"
	"github.com/globalsign/mgo"
	e "../http"
	"../model"
)

type LogsDb struct {
	Collection *mgo.Collection
}

func (d *LogsDb) SelectMessages(writer http.ResponseWriter, query interface{}, pageNr int, pageSizeNr int){
	var result struct {
		Message string
	}

	iter := d.Collection.Find(query).
		Sort("-timestamp").
		Skip(pageNr * pageSizeNr).
		Limit(pageSizeNr).
		Select(bson.M{"message": 1}).Iter()

	for iter.Next(&result) {
		writer.Write([]byte(result.Message + "\n"))
	}
	if iter.Timeout() {
		e.FailOnError(errors.New("Database iterator timeout..."), writer)
	}
	if err := iter.Close(); err != nil {
		e.FailOnError(err, writer)
	}
}

func (d *LogsDb) SelectFullLogs(writer http.ResponseWriter, query interface{}, pageNr int, pageSizeNr int){
	var result []model.LogEntry

	err := d.Collection.Find(query).
		Sort("-timestamp").
		Skip(pageNr * pageSizeNr).
		Limit(pageSizeNr).
		All(&result)

	if err != nil {
		e.FailOnError(err, writer)
		return
	}

	json.NewEncoder(writer).Encode(result)
}

func (d *LogsDb) SelectUiLogs(writer http.ResponseWriter, query interface{}, pageNr int, pageSizeNr int){
	var result struct {
		Id bson.ObjectId `bson:"_id,omitempty"`
		Message string
		Timestamp string
		TimestampSeconds uint32
		Fields map[string]string
		Container struct {Name string}
	}

	iter := d.Collection.Find(query).
		Select(bson.M{"_id": 1, "message": 1, "timestamp": 1, "timestampseconds": 1, "fields": 1, "container.name": 1}).
		Sort("-timestamp").
		Skip(pageNr * pageSizeNr).
		Limit(pageSizeNr).
		Iter()

	encoder := json.NewEncoder(writer)

	writer.Write([]byte("["))
	for iter.Next(&result) {
		encoder.Encode(result)
		if !iter.Done() {
			writer.Write([]byte(","))
		}
	}
	writer.Write([]byte("]"))

	if iter.Timeout() {
		e.FailOnError(errors.New("database iterator timeout"), writer)
	}
	if err := iter.Close(); err != nil {
		e.FailOnError(err, writer)
	}
}