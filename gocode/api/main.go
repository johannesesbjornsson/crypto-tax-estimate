package main

import (
	"github.com/gorilla/mux"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/db"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	db, sqlDB := db.InitDB()
	defer sqlDB.Close()

	r.HandleFunc("/v1/user", func(w http.ResponseWriter, r *http.Request) {
		GetUser(db, w, r)
	}).Methods("GET")
	r.HandleFunc("/v1/user", func(w http.ResponseWriter, r *http.Request) {
		CreateOrUpdateUser(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		GetTransactions(db, w, r)
	}).Methods("GET")

	r.HandleFunc("/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		CreateOrUpdateTransaction(db, w, r)
	}).Methods("POST")
	r.HandleFunc("/v1/transactions/upload", func(w http.ResponseWriter, r *http.Request) {
		UploadCSV(db, w, r)
	}).Methods("POST")
	r.HandleFunc("/v1/transactions/upload", func(w http.ResponseWriter, r *http.Request) {
		GetFileUploads(db, w, r)
	}).Methods("GET")

	log.Infof("Listening on port 8080")
	http.ListenAndServe(":8080", r)

}
