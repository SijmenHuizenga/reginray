package main

import (
	"github.com/globalsign/mgo"
	"log"
	"os"
	"net/http"
	"strings"
	"./decoder"
	"./config"
	"./db"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		log.Fatal("PORT env var not provided")
	}

	mongoServer := os.Getenv("MONGO_SERVER")
	if strings.TrimSpace(mongoServer) == "" {
		log.Fatal("MONGO_SERVER env var not provided")
	}

	servicePatterns, grokPatterns := config.LoadConfig()

	session, err := mgo.Dial(mongoServer)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	mongodb := session.DB("logs")
	mongocol :=  mongodb.C("logs")

	aspicio := Aspicio{
		port:           port,
		decoder:        decoder.NewDecoder(servicePatterns, grokPatterns),
		logsdb:         db.LogsDb{Collection: mongocol},
	}
	aspicio.start()
}

type Aspicio struct {
	port           string
	mongoServer    string
	decoder        decoder.Decoder
	logsdb         db.LogsDb
	logscollection *mgo.Collection
}

func (a *Aspicio) start() {
	router := mux.NewRouter()
	router.HandleFunc("/logs", a.GetLogs).Methods("GET")
	router.HandleFunc("/logs", a.PostLogs).Methods("POST")
	router.HandleFunc("/stats", a.SelectStats).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	log.Println("Starting server on :" + a.port)

	log.Fatal(http.ListenAndServe(":"+a.port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
