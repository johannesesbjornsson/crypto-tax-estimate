package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"io"
	"fmt"
	"encoding/csv"

	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
)

type BinanceParser struct{}

var reAmountAndAsset = regexp.MustCompile(`^([0-9.]+)([A-Z]+)$`)

func parseAmountAndAsset(input string) (float64, string) {
	m := reAmountAndAsset.FindStringSubmatch(input)
	if len(m) != 2 && len(m) != 3 {
		log.Println("Failed to parse amount and asset from:", input)
		return 0, ""
	}
	val, _ := strconv.ParseFloat(m[1], 64)
	return val, m[2]
}

func (b BinanceParser) HeadersMatch(h []string) bool {
	return strings.EqualFold(cleanHeader(h[0]), "Date(UTC)") &&
    strings.EqualFold(cleanHeader(h[1]),   "Pair") &&
    strings.EqualFold(cleanHeader(h[2]),   "Side") &&
    strings.EqualFold(cleanHeader(h[3]),   "Price") &&
    strings.EqualFold(cleanHeader(h[4]),   "Executed") &&
    strings.EqualFold(cleanHeader(h[5]),   "Amount") &&
    strings.EqualFold(cleanHeader(h[6]),   "Fee")
}


func (b BinanceParser) ParseRecord(r []string) (models.Transaction, error) {
	date, err := time.Parse("2006-01-02 15:04:05", r[0])
	if err != nil {
		return models.Transaction{}, err
	}
	
	amount, asset := parseAmountAndAsset(r[4])
	_, quoteAsset := parseAmountAndAsset(r[5])

	price, _ := strconv.ParseFloat(r[3], 64)


	return models.Transaction{
			Date:        date,
			Description: "",
			Type:        r[2],
			Amount:      amount,
			Price:       price,
			Asset:       asset,
			QuoteCurrency: quoteAsset,
			Source:      "CSV Upload",
			UserID:      1,
	}, nil

}


func (b BinanceParser) ParseFile(reader *csv.Reader) ([]models.Transaction, error) {
	var txs []models.Transaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		tx, err := b.ParseRecord(record)
		if err != nil {
			log.Printf("Skipping row: %v", err)
			continue
		}
		fmt.Printf("Parsed: %+v\n", tx)
		txs = append(txs, tx)
	}
	return txs, nil
}





