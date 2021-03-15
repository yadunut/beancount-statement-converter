package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	file, err := os.Open("test.csv")
	handleError(err)
	result, err := parseRevolutTransaction(file)
	handleError(err)
	fmt.Printf(result)
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
	return fmt.Sprintf(formatStr, t.Date, t.Payee, t.Amount, bank, t.Account)
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
		transaction := Transaction{
			Payee:   record[2],
			Account: record[9],
		}
		if i == 0 {
			continue
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

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
