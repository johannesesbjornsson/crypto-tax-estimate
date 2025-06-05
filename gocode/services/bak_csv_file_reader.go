package main

import (
	"fmt"
	"os"
	"time"

	"encoding/csv"
	//"github.com/gocarina/gocsv"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
	"io"
)

type BinanceTransaction struct {
	Date     time.Time `csv:"Date(UTC)"`
	Pair     string    `csv:"Pair"`
	Side     string    `csv:"Side"`
	Price    float64   `csv:"Price"`
	Executed string    `csv:"Executed"`
	Amount   string    `csv:"Amount"`
	Fee      string    `csv:"Fee"`
}

func main() {
	filePath := "/Users/johannesesbjornsson/Downloads/binance_2021-2022.csv"

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		return r // Allows use dot as delimiter and use quotes in CSV
	})
	var transactions []models.Transaction
	err = gocsv.UnmarshalToCallback(file, func(tx *BinanceTransaction) error {
		// Your calculations go here
		fmt.Printf("Processing transaction: %+v\n", tx)

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

		return nil
	})

	if err != nil {
		fmt.Printf("CSV parse error: %v\n", err)
	}

	var binanceTransaction []*BinanceTransaction
	if err := gocsv.Unmarshal(file, &binanceTransaction); err != nil {
		fmt.Printf("CSV parse error: %v\n", err)
		return
	}

	//for _, tx := range transactions {
	//	fmt.Printf("BinanceTransaction: %+v\n", tx)
	//}
}
