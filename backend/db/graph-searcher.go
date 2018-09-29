package db

import (
	"net/http"
	"github.com/globalsign/mgo/bson"
	"encoding/json"
	"errors"
	e "../http"
)

func (d *LogsDb) SelectStats(writer http.ResponseWriter, secondInterval int, rangeAfter int64, rangeEnd int64) {
	// map[_id:map[count:1 _id:map[timestamp:1.53762321e+08]]]
	var result bson.M

	type ResultOutputStruct struct {
		Count     int
		TimestampPerUnit int64
		TimestampSeconds int64
	}

	pipe := d.Collection.Pipe([]bson.M{
		{"$match": bson.M{"timestampseconds": bson.M{"$gt": rangeAfter, "$lt": rangeEnd}}},
		{"$group": bson.M{
			"_id":   bson.M{"timestamp": bson.M{"$trunc": bson.M{"$divide": []interface{}{"$timestampseconds", secondInterval}}}},
			"count": bson.M{"$sum": 1},
		}},
	})

	iter := pipe.Iter()

	encoder := json.NewEncoder(writer)

	writer.Write([]byte("["))
	for iter.Next(&result) {
		y := int64(result["_id"].(bson.M)["timestamp"].(float64))
		encoder.Encode(ResultOutputStruct{
			Count:     result["count"].(int),
			TimestampPerUnit: y,
			TimestampSeconds: y * int64(secondInterval),
		})
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
