package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

type Transaction struct {
	Date time.Time
}

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

func parseRevolutTransaction(file io.Reader) (string,error) {
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	records, err := csvReader.ReadAll()
	if err != nil {
		return "",err
	}

	results := make([]string, 20)

	for i, record := range records {
		// Skip first line in CSV as headers
		if i == 0 {
			continue
		}
		time, err := time.Parse("2 Jan 2006", record[0])
		if err != nil {
			return "", err
		}
		timeStr := time.Format("2006-01-02")
		payee := record[2]
		account := record[9]

		amount := ""
		if record[3] != "" {
			amount = "-" + record[3]
		} else if record[4] != "" {
			amount = record[4]
		} else {
			return "", fmt.Errorf("no paid in / paid out amount")
		}
		formatStr :=`%s ! "%s" ""
  Assets:Revolut:SGD %s SGD
  Expenses:%s

`
		results = append(results, fmt.Sprintf(formatStr , timeStr, payee, amount, account))
	}
	reverse(results)
	result := ""
	for _, line := range results  {
		result += line
	}

	return result, nil
}

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}
