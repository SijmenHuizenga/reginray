package main

import (
	"net/http"
	"strconv"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"errors"
	e "./http"
)

func (a *Aspicio) GetLogs(writer http.ResponseWriter, request *http.Request) {
	queries := request.URL.Query()

	var findQuery = bson.M{}

	////////////////////////////////////
	period := queries.Get("timeperiod")
	if len(period) != 0 {
		after, err1 := strconv.ParseInt(strings.Split(period, "-")[0], 10, 64)
		before, err2 := strconv.ParseInt(strings.Split(period, "-")[1], 10, 64)
		if err1 != nil || err2 != nil {
			e.FailOnBadRequest(errors.New("timeperiod not a number format of pattern: 123-456"), writer)
			return
		}
		findQuery["timestampseconds"] = bson.M{"$gt": after, "$lt": before}
	} else {
		e.FailOnBadRequest(errors.New("timeperiod get argument is required"), writer)
		return
	}

	////////////////////////////////////
	page := queries.Get("page")
	if len(page) == 0 {
		e.FailOnBadRequest(errors.New("page get argument is required"), writer)
		return
	}
	pageNr, err := strconv.Atoi(page)
	if err != nil {
		e.FailOnBadRequest(errors.New("page must be a number"), writer)
		return
	}
	/////////////////////////////////////////
	pageSize := queries.Get("pagesize")
	if len(pageSize) == 0 {
		e.FailOnBadRequest(errors.New("pagesize get argument is required"), writer)
		return
	}
	pageSizeNr, err := strconv.Atoi(pageSize)
	if err != nil {
		e.FailOnBadRequest(errors.New("pagesize must be a number"), writer)
		return
	}
	/////////////////////////////////////////

	var format = queries.Get("format")
	if len(format) == 0 {
		format = "ui"
	}

	////////////////////////////////////////

	switch format {
	case "full":
		a.logsdb.SelectFullLogs(writer, findQuery, pageNr, pageSizeNr)
	case "messages":
		a.logsdb.SelectMessages(writer, findQuery, pageNr, pageSizeNr)
	case "ui":
		a.logsdb.SelectUiLogs(writer, findQuery, pageNr, pageSizeNr)
	}
}

func (a *Aspicio) SelectStats(writer http.ResponseWriter, request *http.Request) {
	queries := request.URL.Query()

	interval := queries.Get("interval")
	if len(interval) == 0 {
		e.FailOnBadRequest(errors.New("interval get argument is required"), writer)
		return
	}
	intervalNr, err3 := strconv.Atoi(interval)
	if err3 != nil {
		e.FailOnBadRequest(errors.New("interval must be a number"), writer)
		return
	}


	period := queries.Get("timeperiod")
	if len(period) == 0 {
		e.FailOnBadRequest(errors.New("timeperiod get argument is required"), writer)
		return
	}

	after, err1 := strconv.ParseInt(strings.Split(period, "-")[0], 10, 64)
	before, err2 := strconv.ParseInt(strings.Split(period, "-")[1], 10, 64)
	if err1 != nil || err2 != nil {
		e.FailOnBadRequest(errors.New("timeperiod not a number format of pattern: 123-456"), writer)
		return
	}
	a.logsdb.SelectStats(writer, intervalNr, after, before)
}