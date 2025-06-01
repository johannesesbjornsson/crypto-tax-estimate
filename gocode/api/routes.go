package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/johannesesbjornsson/crypto-tax-estimate/database/db"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
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

	transactions, err := db.GetTransactionsByEmail(email, limit, offset)
	if err != nil {
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func CreateOrUpdateTransaction(db *db.Database, w http.ResponseWriter, r *http.Request) {
	var tx models.Transaction
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

	if err := db.CreateTransaction(&tx); err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
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

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	description := r.FormValue("description")

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	_, err = reader.Read() // Skip header
	if err != nil {
		http.Error(w, "Invalid CSV header", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByEmail("johannes.esbjornsson@gmail.com")
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

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

	var transactions []models.Transaction

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 7 {
			continue
		}

		date, err := time.Parse("2006-01-02 15:04:05", record[0])
		if err != nil {
			log.Warnf("Invalid date format: %v", record[0])
			continue
		}

		amountField := record[4]
		var amount float64
		var asset string
		var amountAssetRegexp = regexp.MustCompile(`^([0-9.]+)([A-Za-z]+)$`)

		// Inside your loop
		matches := amountAssetRegexp.FindStringSubmatch(amountField)
		if len(matches) != 3 {
			log.Warnf("Failed to parse amount and asset from: %q", amountField)
			continue
		}
		amount, err = strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Warnf("Invalid amount value: %q", matches[1])
			continue
		}
		asset = matches[2]

		price, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Warnf("Invalid price: %v", record[3])
			continue
		}

		tx := models.Transaction{
			Date:        date,
			Description: description,
			Type:        strings.Title(strings.ToLower(record[2])),
			Amount:      amount,
			Price:       price,
			Asset:       asset,
			Source:      "CSV Upload",
			UserID:      user.ID,
		}

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
