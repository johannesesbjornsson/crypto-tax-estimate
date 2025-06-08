package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/johannesesbjornsson/crypto-tax-estimate/database/db"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
	log "github.com/sirupsen/logrus"
)

const (
	limitHours = 24 // Number of hours to fetch
	apiURL     = "https://min-api.cryptocompare.com/data/v2/histohour"
	quoteUSD   = "USD"
)

type CryptoCompareResponse struct {
	Data struct {
		Data []struct {
			Time  int64   `json:"time"`
			Close float64 `json:"close"`
		} `json:"Data"`
	} `json:"Data"`
}

func fetchHourlyPrices(symbol string) ([]models.MarketPrice, error) {
	url := fmt.Sprintf("%s?fsym=%s&tsym=%s&limit=%d", apiURL, symbol, quoteUSD, limitHours)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request price data: %w", err)
	}
	defer resp.Body.Close()

	var result CryptoCompareResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var prices []models.MarketPrice
	for _, entry := range result.Data.Data {
		prices = append(prices, models.MarketPrice{
			Price:       entry.Close,
			Timestamp:   time.Unix(entry.Time, 0),
			QuoteCurrencyID: 1, // assume USD has ID = 1
		})
	}
	return prices, nil
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	database, sqlDB := db.InitDB()
	defer sqlDB.Close()

	currencies, err := database.GetCurrencies()
	if err != nil {
		log.Fatalf("Failed to fetch currencies: %v", err)
	}

	var usdID uint
	for _, c := range currencies {
		if c.Name == quoteUSD {
			usdID = c.ID
			break
		}
	}
	if usdID == 0 {
		log.Fatalf("USD currency not found in database")
	}

	for _, currency := range currencies {
		if currency.Name == quoteUSD {
			continue
		}

		log.Infof("Fetching price data for %s -> USD", currency.Name)
		prices, err := fetchHourlyPrices(currency.Name)
		if err != nil {
			log.Errorf("Failed to fetch price data for %s: %v", currency.Name, err)
			continue
		}

		for _, price := range prices {
			price.BaseCurrencyID = currency.ID
			price.QuoteCurrencyID = usdID
			if err := database.DB.Clauses(
				// Upsert
				// "ON CONFLICT (base_currency_id, quote_currency_id, timestamp) DO NOTHING"
			).Create(&price).Error; err != nil {
				log.Errorf("Failed to insert price for %s at %v: %v", currency.Name, price.Timestamp, err)
			}
		}
	}
}