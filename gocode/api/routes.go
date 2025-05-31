package main

import (
	"encoding/json"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/db"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetUser(db *db.Database, w http.ResponseWriter, r *http.Request) {
	log.Infof("Received request: %s %s", r.Method, r.URL.Path)

	email := "johannes.esbjornsson@gmail.com"

	user, err := db.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateOrUpdateUser(db *db.Database, w http.ResponseWriter, r *http.Request) {
	log.Infof("Received request: %s %s", r.Method, r.URL.Path)

	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := db.CreateOrUpdateUser(&input); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(input)
}


func GetTransactions(db *db.Database, w http.ResponseWriter, r *http.Request) {
	log.Infof("Received request: %s %s", r.Method, r.URL.Path)

  email := "johannes.esbjornsson@gmail.com"

    transactions, err := db.GetTransactionsByEmail(email)
    if err != nil {
        log.Errorf("Error fetching transactions: %v", err)
        http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)

}