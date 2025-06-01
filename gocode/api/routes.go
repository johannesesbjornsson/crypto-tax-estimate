package main

import (
	"encoding/json"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/db"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"encoding/csv"
	"io"
	"strconv"
	"fmt"
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


func CreateOrUpdateTransaction(db *db.Database, w http.ResponseWriter, r *http.Request) {
	log.Infof("Received request: %s %s", r.Method, r.URL.Path)

	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		log.Errorf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hardcoded user email, as per GetTransactions
	email := "johannes.esbjornsson@gmail.com"
	user, err := db.GetUserByEmail(email)
	if err != nil {
		log.Errorf("Failed to find user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	if tx.Date.After(today) {
		http.Error(w, "Transaction date cannot be in the future", http.StatusBadRequest)
		return
	}

	tx.UserID = user.ID
	tx.Source = "Manual"

	if err := db.CreateTransaction(&tx); err != nil {
		log.Errorf("Failed to create transaction: %v", err)
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

func UploadCSV(db *db.Database, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header row
	if _, err := reader.Read(); err != nil {
		http.Error(w, "Invalid CSV header", http.StatusBadRequest)
		return
	}

	var transactions []models.Transaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 7 {
			continue
		}

		date, _ := time.Parse("2006-01-02", record[0])
		amount, _ := strconv.ParseFloat(record[4], 64)

		tx := models.Transaction{
			Date:        date,
			Description: record[1],
			Venue:       record[2],
			Type:        record[3],
			Amount:      amount,
			Asset:       record[5],
			Source:      record[6],
		}

		user, err := db.GetUserByEmail("johannes.esbjornsson@gmail.com") // or derive dynamically
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		tx.UserID = user.ID

		transactions = append(transactions, tx)
	}

	for _, tx := range transactions {
		if err := db.CreateTransaction(&tx); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("CSV upload successful"))

}