package main

import (
	//"encoding/csv"
	"encoding/json"
	"fmt"
	//"io"
	"net/http"
	//"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/johannesesbjornsson/crypto-tax-estimate/database/db"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
	csv_parser "github.com/johannesesbjornsson/crypto-tax-estimate/services/csv-parser"
	log "github.com/sirupsen/logrus"
)

func GetUser(db *db.Database, w http.ResponseWriter, r *http.Request) {
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
	email := "johannes.esbjornsson@gmail.com"

	// Parse query parameters
	limit := 100
	offset := 0
	txType := "trade"

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	if o := r.URL.Query().Get("txType"); o != "" {
		txType = o
	}
	var response interface{}
	if txType == "trade" {
		transactions, totalPages, err := db.GetTradeTransactionsByEmail(email, limit, offset)
		if err != nil {
			http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
			return
		}
		response = struct {
			Transactions []models.TradeTransaction `json:"transactions"`
			TotalPages   int                       `json:"totalPages"`
		}{
			Transactions: transactions,
			TotalPages:   totalPages,
		}
	} else if txType == "simple" {
		transactions, totalPages, err := db.GetSimpleTransactionsByEmail(email, limit, offset)
		if err != nil {
			http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
			return
		}

		response = struct {
			Transactions []models.SimpleTransaction `json:"transactions"`
			TotalPages   int                        `json:"totalPages"`
		}{
			Transactions: transactions,
			TotalPages:   totalPages,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateOrUpdateTransaction(db *db.Database, w http.ResponseWriter, r *http.Request) {
	var tx models.BaseTransaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		log.Errorf("Failed to decode transaction JSON: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	email := "johannes.esbjornsson@gmail.com"
	user, err := db.GetUserByEmail(email)
	if err != nil {
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
	tx.Type = strings.ToLower(tx.Type)

	if tx.Type == "buy" || tx.Type == "sell" {
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		tradeTx, _ := db.NewTradeTransaction(&tx, price, r.FormValue("quote_currency"))
		if err := db.CreateTradeTransaction(tradeTx); err != nil {
			http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
			return
		}
	} else {
		simpleTx, _ := db.NewSimpleTransaction(&tx)
		if err := db.CreateSimpleTransaction(simpleTx); err != nil {
			http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

func GetFileUploads(db *db.Database, w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserByEmail("johannes.esbjornsson@gmail.com") // Or use auth context
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	uploads, err := db.GetFileUploadsByUserID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploads)
}

func UploadCSV(db *db.Database, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByEmail("johannes.esbjornsson@gmail.com")
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	description := r.FormValue("description")
	file, fileHeader, err := r.FormFile("file")
	simpleTransactions, tradeTransactions, err := csv_parser.ParseCSV(file)
	if err != nil {
		log.Errorf("Failed to parse CSV file: %v", err)
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	log.Infof("Parsed %d trade transactions from CSV file %s", len(tradeTransactions)+len(simpleTransactions), fileHeader.Filename)

	fileUpload := models.FileUploads{
		Name:        fileHeader.Filename,
		UserID:      user.ID,
		Description: description,
		CreatedAt:   time.Now(),
	}

	if err := db.CreateFileUpload(&fileUpload); err != nil {
		http.Error(w, "Failed to record file upload", http.StatusInternalServerError)
		return
	}

	for _, tx := range tradeTransactions {
		tx.Description = description
		tx.Source = fileHeader.Filename
		if err := db.CreateTradeTransaction(&tx); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save: %v", err), http.StatusInternalServerError)
			return
		}
	}
	for _, tx := range simpleTransactions {
		tx.Description = description
		tx.Source = fileHeader.Filename
		if err := db.CreateSimpleTransaction(&tx); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("CSV upload successful"))

}
