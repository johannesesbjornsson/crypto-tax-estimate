package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	r.HandleFunc("/v1/user", GetUser).Methods("GET")
	r.HandleFunc("/v1/user", CreateOrUpdateUser).Methods("POST")

	log.Infof("Listening on port 8080")
	http.ListenAndServe(":8080", r)

}
