package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
	log "github.com/sirupsen/logrus"
)

type KrakenParser struct{}

type krakenGenericRecord struct {
	Txid    string    `json:"txid"`
	Refid   string    `json:"refid"`
	Time    time.Time `json:"time"`
	Type    string    `json:"type"`
	Subtype string    `json:"subtype"`
	Aclass  string    `json:"aclass"`
	Asset   string    `json:"asset"`
	Wallet  string    `json:"wallet"`
	Amount  float64   `json:"amount"`
	Fee     float64   `json:"fee"`
	Balance float64   `json:"balance"`
}

var reStakedAssetName = regexp.MustCompile(`^([A-Z]+)([0-9.]*\.S)?$`)

func parseStakedAssetName(input string) (string, error) {
	m := reStakedAssetName.FindStringSubmatch(input)
	return m[1], nil

}

func isBaseCurrency(asset string) bool {
	baseCurrencies := []string{"GBP", "EUR", "USD", "JPY", "USDC", "USDT"}
	for _, base := range baseCurrencies {
		if strings.EqualFold(asset, base) {
			return true
		}
	}
	return false
}

func (k KrakenParser) HeadersMatch(h []string) bool {
	return strings.EqualFold(cleanHeader(h[0]), "txid") &&
		strings.EqualFold(cleanHeader(h[1]), "refid") &&
		strings.EqualFold(cleanHeader(h[2]), "time") &&
		strings.EqualFold(cleanHeader(h[3]), "type") &&
		strings.EqualFold(cleanHeader(h[4]), "subtype") &&
		strings.EqualFold(cleanHeader(h[5]), "aclass") &&
		strings.EqualFold(cleanHeader(h[6]), "asset") &&
		strings.EqualFold(cleanHeader(h[7]), "wallet") &&
		strings.EqualFold(cleanHeader(h[8]), "amount") &&
		strings.EqualFold(cleanHeader(h[9]), "fee") &&
		strings.EqualFold(cleanHeader(h[10]), "balance")
}

func (k KrakenParser) ParseGenericRecord(r []string) (krakenGenericRecord, error) {
	date, err := time.Parse("2006-01-02 15:04:05", r[2])
	if err != nil {
		return krakenGenericRecord{}, err
	}
	amount, _ := strconv.ParseFloat(r[8], 64)
	fee, _ := strconv.ParseFloat(r[9], 64)
	balance, _ := strconv.ParseFloat(r[10], 64)

	return krakenGenericRecord{
		Txid:    r[0],
		Refid:   r[1],
		Time:    date,
		Type:    r[3],
		Subtype: r[4],
		Aclass:  r[5],
		Asset:   r[6],
		Wallet:  r[7],
		Amount:  amount,
		Fee:     fee,
		Balance: balance,
	}, nil
}

func (k KrakenParser) ParseTradeRecord(recieve krakenGenericRecord, taken krakenGenericRecord) (models.TradeTransaction, error) {
	if recieve.Refid != taken.Refid {
		return models.TradeTransaction{}, fmt.Errorf("both records must be of same ref ix ", recieve.Refid, taken.Refid)
	}

	isBaseCurrency := isBaseCurrency(taken.Asset)
	side := "buy"
	if !isBaseCurrency {
		side = "sell"
	}
	if taken.Amount > 0 {
		return models.TradeTransaction{}, fmt.Errorf("taken amount must be negative", taken.Amount)
	}

	return models.TradeTransaction{
		QuoteCurrency: taken.Asset,
		Type:          side,
		Price:         math.Abs(taken.Amount) / recieve.Amount,
		BaseTransaction: models.BaseTransaction{
			Date:        recieve.Time,
			ExternalID:  taken.Txid,
			Description: "",
			Amount:      recieve.Amount,
			Asset:       recieve.Asset,
			Source:      "CSV Upload",
			UserID:      1,
		},
	}, nil
}

func (k KrakenParser) ParseSimpleRecord(record krakenGenericRecord) (models.SimpleTransaction, error) {

	assetName, err := parseStakedAssetName(record.Asset)
	if err != nil {
		return models.SimpleTransaction{}, fmt.Errorf("failed to parse asset name %s: %v", record.Asset, err)
	}
	return models.SimpleTransaction{
		Type: "income",
		BaseTransaction: models.BaseTransaction{
			Date:        record.Time,
			ExternalID:  record.Txid,
			Description: "",
			Amount:      record.Amount,
			Asset:       assetName,
			Source:      "CSV Upload",
			UserID:      1,
		},
	}, nil
}

func (b KrakenParser) ParseFile(reader *csv.Reader) ([]models.SimpleTransaction, []models.TradeTransaction, error) {
	var tradeTxs []models.TradeTransaction
	var simpleTxs []models.SimpleTransaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		genericRecord, err := b.ParseGenericRecord(record)
		if err != nil {
			log.Infof("Skipping row: %v", err)
			continue
		}

		if genericRecord.Type == "trade" {
			relatedRecord, err := reader.Read()
			if err == io.EOF {
				log.Infof("No related record found for trade, skipping")
				break
			}
			relatedGenericRecord, err := b.ParseGenericRecord(relatedRecord)
			tx, err := b.ParseTradeRecord(relatedGenericRecord, genericRecord)
			if err != nil {
				log.Infof("Skipping row: %v", err)
				continue
			}
			tradeTxs = append(tradeTxs, tx)
		} else if genericRecord.Type == "staking" || genericRecord.Type == "earn" {
			if genericRecord.Subtype == "migration" {
				log.Infof("Skipping migration record: %s", genericRecord.Txid)
				continue
			}
			tx, err := b.ParseSimpleRecord(genericRecord)
			if err != nil {
				log.Warnf("Skipping row: %v", err)
				continue
			}
			simpleTxs = append(simpleTxs, tx)
		} else {
			continue
		}
	}

	return simpleTxs, tradeTxs, nil
}
