package csvparser

import (
	"encoding/csv"
	"fmt"
	"github.com/johannesesbjornsson/crypto-tax-estimate/database/models"
	"log"
	"mime/multipart"
	"strings"
)

type CSVParser interface {
	HeadersMatch([]string) bool
	ParseFile(*csv.Reader) ([]models.SimpleTransaction, []models.TradeTransaction, error)
}

func cleanHeader(s string) string {
	if strings.HasPrefix(s, "\uFEFF") {
		s = strings.TrimPrefix(s, "\uFEFF")
	}
	return strings.Trim(s, `"`)
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


     
func ParseCSV(file multipart.File) ([]models.SimpleTransaction, []models.TradeTransaction, error) {
	//filePath := "/Users/johannesesbjornsson/workspace/personal-testing/binance_2021-2022.csv"
	//filePath := "/Users/johannesesbjornsson/workspace/personal-testing/kraken_2023-2024.csv"

	//file, err := os.Open(filePath)
	//if err != nil {
	//	log.Fatalf("Failed to open file: %v", err)
	//}
	//defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	//reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read header: %v", err)
	}

	parsers := []CSVParser{BinanceParser{}, KrakenParser{}}
	parser, err := detectParser(headers, parsers)
	if err != nil {
		log.Fatalf("Could not detect CSV format: %v", err)
	}

	fmt.Printf("Detected format: %T\n", parser)

	return parser.ParseFile(reader)

}

/*
func main() {
	filePath := "/Users/johannesesbjornsson/workspace/personal-testing/binance_2021-2022.csv"
	//filePath := "/Users/johannesesbjornsson/workspace/personal-testing/kraken_2023-2024.csv"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()
	ParseFile(file) // Call the ParseFile function to start parsing

}
*/