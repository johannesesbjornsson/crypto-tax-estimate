package main

import (
	"strconv"
	"strings"
	"time"
	"encoding/csv"

	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
)


type KrakenParser struct{}

func (k KrakenParser) HeadersMatch(h []string) bool {
	return len(h) > 9 && strings.EqualFold(h[0], "txid") && strings.EqualFold(h[2], "time")
}

func (k KrakenParser) ParseRecord(r []string) (models.Transaction, error) {
	ts, err := strconv.ParseFloat(r[2], 64)
	if err != nil {
		return models.Transaction{}, err
	}
	date := time.Unix(int64(ts), 0)
	amount, _ := strconv.ParseFloat(r[8], 64)

	return models.Transaction{
		Date:     date,
		Description: "",
		Type:     r[3],
		Asset:    r[6],
		Amount:   amount,
		Price:    0,
		QuoteCurrency: "USD",
		Source:   "Kraken",
	}, nil
}



func (b KrakenParser) ParseFile(reader *csv.Reader) ([]models.Transaction, error) {
	var txs []models.Transaction
	return txs, nil
}