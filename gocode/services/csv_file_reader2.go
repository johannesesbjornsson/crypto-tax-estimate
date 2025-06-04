package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
)


type CSVParser interface {
	HeadersMatch([]string) bool
	ParseRecord([]string) (models.Transaction, error)
}

func cleanHeader(s string) string {
	if strings.HasPrefix(s, "\uFEFF") {
		s = strings.TrimPrefix(s, "\uFEFF")
	}
  return strings.Trim(s, `"`)
}

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

//func parsePair(input string) (string, string) {
//	knownCurrencies := []string{
//		"ADA" , "BUSD", "XRP", "USDC", "XLM", "SOL", "SHI", "ROSE", "MATIC", "BNB", "FTM", "DOT", "GBP",
//	}
//	return "", ""
//
//}
// ------------------------------------
// Binance Parser
// ------------------------------------

type BinanceParser struct{}

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

// ------------------------------------
// Coinbase Parser
// ------------------------------------

type CoinbaseParser struct{}

func (c CoinbaseParser) HeadersMatch(h []string) bool {
	return len(h) > 6 && strings.EqualFold(h[0], "ID") && strings.Contains(h[1], "Timestamp")
}

func (c CoinbaseParser) ParseRecord(r []string) (models.Transaction, error) {
	date, err := time.Parse("2006-01-02T15:04:05Z", r[1])
	if err != nil {
		return models.Transaction{}, err
	}
	amount, _ := strconv.ParseFloat(r[4], 64)
	price, _ := strconv.ParseFloat(r[6], 64)

	return models.Transaction{
			Date:        date,
			Description: "",
			Type:        r[2],
			Amount:      amount,
			Price:       price,
			Asset:       r[3],
			QuoteCurrency: "USD",
			Source:      "CSV Upload",
			UserID:      1,
	}, nil
}

// ------------------------------------
// Kraken Parser
// ------------------------------------

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



func detectParser(headers []string, parsers []CSVParser) (CSVParser, error) {
	for _, p := range parsers {
		fmt.Printf("Checking parser: %T\n", p)
		if p.HeadersMatch(headers) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("unrecognized CSV format")
}


func main() {
	filePath := "/Users/johannesesbjornsson/Downloads/binance_2021-2022.csv"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read header: %v", err)
	}

	parsers := []CSVParser{BinanceParser{}, CoinbaseParser{}, KrakenParser{}}
	parser, err := detectParser(headers, parsers)
	if err != nil {
		log.Fatalf("Could not detect CSV format: %v", err)
	}

	fmt.Printf("Detected format: %T\n", parser)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < len(headers) {
			continue
		}
		tx, err := parser.ParseRecord(record)
		if err != nil {
			log.Printf("Skipping row: %v", err)
			continue
		}
		fmt.Printf("Parsed: %+v\n", tx)
	}
}