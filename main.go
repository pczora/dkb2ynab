package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pczora/dkb2ynab/config"
)

func main() {
	args := os.Args
	config := config.NewDkbConfig()
	f, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bufferedReader := bufio.NewReader(f)
	for i := 0; i < config.SkipLines; i++ {
		bufferedReader.ReadBytes('\n')
	}
	reader := csv.NewReader(bufferedReader)
	reader.Comma = ';'
	csvRecords, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	var parsedRecords []Record
	for _, r := range csvRecords {
		date, err := time.Parse(config.DateFormat, r[config.DateColumn])
		if err != nil {
			panic(err)
		}
		amountString := r[config.AmountColumn]
		if config.NormalizeAmount {
			normalizeAmount(&amountString)
		}
		amount, err := strconv.ParseFloat(amountString, 64)
		if err != nil {
			panic(err)
		}
		record := Record{date, r[config.PayeeColumn], r[config.MemoColumn], amount}
		parsedRecords = append(parsedRecords, record)
	}
	writeRecords(parsedRecords)
}

// Attention: this currently only normalizes amounts in German formatting
func normalizeAmount(amount *string) {
	*amount = strings.Replace(*amount, ".", "", -1)
	*amount = strings.Replace(*amount, ",", ".", -1)
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
