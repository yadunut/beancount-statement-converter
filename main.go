package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"flag"
	"time"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	fileName string
	bank string
)

func init() {
	flag.StringVar(&fileName,"f", "statement.csv", "filename of csv file to parse")
	flag.StringVar(&bank,"bank", "", "dbs or revolut")
}

func main() {
	flag.Parse()
	file, err := os.Open(fileName)
	handleError(err)
	switch bank {
	case "dbs": {
		result, err := parseDbsTransaction(file)
		handleError(err)
		fmt.Printf(result)
	}
	case "revolut": {
		result, err := parseRevolutTransaction(file)
		handleError(err)
		fmt.Printf(result)
	}
	default: {
		handleError(fmt.Errorf("invalid bank: %s", bank))
	}
	}
}

type Transaction struct {
	Date    string
	Payee   string
	Account string
	Amount  string
}

func (t *Transaction) String(bank string) string {
	formatStr := `%s ! "%s" ""
  Assets:%s:SGD %s SGD
  Expenses:%s

`
	return fmt.Sprintf(formatStr, t.Date, t.Payee, bank, t.Amount, t.Account)
}

func parseRevolutTransaction(file io.Reader) (string, error) {
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	records, err := csvReader.ReadAll()
	if err != nil {
		return "", err
	}

	results := make([]string, 20)

	for i, record := range records {
		// Skip first line in CSV as headers
		if i == 0 {
			continue
		}

		transaction := Transaction{
			Payee:   record[2],
			Account: record[9],
		}
		time, err := time.Parse("2 Jan 2006", record[0])
		if err != nil {
			return "", err
		}
		transaction.Date = time.Format("2006-01-02")

		amount := ""
		if record[3] != "" {
			amount = "-" + record[3]
		} else if record[4] != "" {
			amount = record[4]
		} else {
			return "", fmt.Errorf("no paid in / paid out amount")
		}
		transaction.Amount = amount
		results = append(results, transaction.String("Revolut"))
	}

	reverse(results)
	result := ""
	for _, line := range results {
		result += line
	}

	return result, nil
}

func parseDbsTransaction(file io.Reader) (string, error) {
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return "", err
	}

	results := make([]string, 20)

	for i, record := range records {
		// Skip first line in CSV as headers
		transaction := Transaction{
			Payee:   record[4],
			Account: "",
		}
		if i == 0 {
			continue
		}
		time, err := time.Parse("02 Jan 2006", record[0])
		if err != nil {
			return "", err
		}

		transaction.Date = time.Format("2006-01-02")

		amount := ""
		if strings.TrimSpace(record[2]) != "" {
			amount = "-" + strings.TrimSpace(record[2])
		} else if strings.TrimSpace(record[3]) != "" {
			amount = strings.TrimSpace(record[3])
		} else {
			return "", fmt.Errorf("no paid in / paid out amount")
		}
		transaction.Amount = amount
		results = append(results, transaction.String("DBS"))
	}


	reverse(results)
	result := ""
	for _, line := range results {
		result += line
	}

	return result, nil

}

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
