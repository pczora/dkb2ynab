package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args
	f, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(f)
	reader.Comma = ';'
	csvRecords, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	var parsedRecords []Record
	for _, r := range csvRecords {
		date, err := time.Parse("02.01.2006", r[1])
		if err != nil {
			panic(err)
		}
		normalizedAmount := normalizeAmount(r[7])
		amount, err := strconv.ParseFloat(normalizedAmount, 64)
		if err != nil {
			panic(err)
		}
		record := Record{date, r[3], r[4], amount}
		parsedRecords = append(parsedRecords, record)
	}
	writeRecords(parsedRecords)
}

func normalizeAmount(amount string) string {
	intermediatResult := strings.Replace(amount, ".", "", -1)
	result := strings.Replace(intermediatResult, ",", ".", -1)
	return result
}

func writeRecords(records []Record) {
	fmt.Println("\"Date\",\"Payee\",\"Memo\",\"Amount\"")
	for _, r := range records {
		fmt.Printf("\"%v\",\"%v\",\"%v\",\"%v\"\n", r.Date.Format("2006/02/01"), r.Payee, r.Memo, r.Amount)
	}
}

type Record struct {
	Date   time.Time
	Payee  string
	Memo   string
	Amount float64
}
