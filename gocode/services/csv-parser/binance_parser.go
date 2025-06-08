package csvparser

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
)

type BinanceParser struct{}

var reAmountAndAsset = regexp.MustCompile(`^([0-9.]+)([A-Z]+)$`)

func parseAmountAndAsset(input string) (float64, string) {
	m := reAmountAndAsset.FindStringSubmatch(input)
	if len(m) != 2 && len(m) != 3 {
		log.Infof("Failed to parse amount and asset from:", input)
		return 0, ""
	}
	val, _ := strconv.ParseFloat(m[1], 64)
	return val, m[2]
}

func (b BinanceParser) HeadersMatch(h []string) bool {
	return strings.EqualFold(cleanHeader(h[0]), "Date(UTC)") &&
		strings.EqualFold(cleanHeader(h[1]), "Pair") &&
		strings.EqualFold(cleanHeader(h[2]), "Side") &&
		strings.EqualFold(cleanHeader(h[3]), "Price") &&
		strings.EqualFold(cleanHeader(h[4]), "Executed") &&
		strings.EqualFold(cleanHeader(h[5]), "Amount") &&
		strings.EqualFold(cleanHeader(h[6]), "Fee")
}

func (b BinanceParser) ParseTradeRecord(r []string) (models.TradeTransaction, error) {
	date, err := time.Parse("2006-01-02 15:04:05", r[0])
	if err != nil {
		return models.TradeTransaction{}, err
	}

	amount, asset := parseAmountAndAsset(r[4])
	_, quoteAsset := parseAmountAndAsset(r[5])

	price, _ := strconv.ParseFloat(r[3], 64)

	return models.TradeTransaction{
		Type:          strings.ToLower(r[2]),
		Price:         price,
		QuoteCurrency: quoteAsset,
		BaseTransaction: models.BaseTransaction{
			Date:        date,
			Description: "",
			Amount:      amount,
			Asset:       asset,
			Source:      "CSV Upload",
			UserID:      1,
		},
	}, nil

}

func (b BinanceParser) ParseFile(reader *csv.Reader) ([]models.SimpleTransaction, []models.TradeTransaction, error) {
	var txs []models.TradeTransaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		tx, err := b.ParseTradeRecord(record)
		if err != nil {
			log.Infof("Skipping row: %v", err)
			continue
		}
		txs = append(txs, tx)
	}
	return []models.SimpleTransaction{}, txs, nil
}
