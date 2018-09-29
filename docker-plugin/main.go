package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/sdk"
	"os"
)

var backendUrl string

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	backendUrl = os.Getenv("ASPICIO_BACKEND")
	if backendUrl == "" {
		panic("ASPICIO_BACKEND not provided")
	}

	h := sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	handlers(&h, newDriver())
	if err := h.ServeUnix("aspicio", 0); err != nil {
		panic(err)
	}
}