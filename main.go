package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/pczora/dkb2ynab/formats"
)

func main() {
	args := os.Args
	var ynabRecords []formats.YnabRecord
	var dkbRecords []formats.DkbRecord
	err := readCSVRecords(args[1], &dkbRecords)
	if err != nil {
		panic(err)
	}
	for _, r := range dkbRecords {
		ynabRecord := formats.YnabRecord{Date: formats.YnabDateTime(r.Date), Payee: r.Payee, Memo: r.PostingText, Amount: formats.YnabAmount(r.Amount)}
		ynabRecords = append(ynabRecords, ynabRecord)
	}
	marshalled, err := gocsv.MarshalString(ynabRecords)
	if err != nil {
		panic(err)
	}
	fmt.Print(marshalled)
}

func readCSVRecords(path string, out *[]formats.DkbRecord) error {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		reader := csv.NewReader(in)
		reader.Comma = ';'
		return reader
	})
	err = gocsv.UnmarshalFile(f, out)
	if err != nil {
		panic(err)
	}
	return nil
}
