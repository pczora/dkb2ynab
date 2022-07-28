package main

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/pczora/dkb2ynab/formats"
)

func main() {
	args := os.Args
	var ynabRecords []formats.YnabRecord
	var dkbRecords []formats.InternalRecord
	dkbConverter := formats.DkbFormatConverter{}
	ynabConverter := formats.YnabFormatConverter{}
	dkbRecords = dkbConverter.ConvertFromFile(args[1])
	for _, r := range dkbRecords {
		ynabRecord := ynabConverter.ConvertFromInternalRecord(r)
		ynabRecords = append(ynabRecords, ynabRecord)
	}
	marshalled, err := gocsv.MarshalString(ynabRecords)

	if err != nil {
		panic(err)
	}
	fmt.Print(marshalled)
}
