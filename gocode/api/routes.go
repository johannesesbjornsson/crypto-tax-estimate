package main

import (
    "encoding/json"
    "net/http"
    log "github.com/sirupsen/logrus"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
		log.Infof("Received request: %s %s", r.Method, r.URL.Path)
    json.NewEncoder(w).Encode(map[string]string{"message": "User endpoint"})
}
