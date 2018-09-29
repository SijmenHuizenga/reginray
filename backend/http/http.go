package http

import (
	"net/http"
	"html"
)

func FailOnError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error": "` + html.EscapeString(err.Error()) + `"}`))
}

func FailOnBadRequest(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"error": "` + html.EscapeString(err.Error()) + `"}`))
}
