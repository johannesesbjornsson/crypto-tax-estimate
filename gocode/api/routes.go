package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/db"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
)

func GetUser(w http.ResponseWriter, r *http.Request) {

	log.Infof("Received request: %s %s", r.Method, r.URL.Path)
	conn := db.InitDB()
	var user models.User

	if err := conn.First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
