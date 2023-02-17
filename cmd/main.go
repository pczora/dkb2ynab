package main

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/pczora/dkb2ynab/formats"
)

func main() {
	args := os.Args
	//TODO: validation
	path := args[1]
	var ynabRecords []formats.YnabRecord
	var inputRecords []formats.InternalRecord
	ynabConverter := formats.YnabFormatConverter{}

	converter, err := formats.FindSuitableConverter(path)
	if err != nil {
		panic(err)
	}
	records, err := converter.ConvertFromFile(path)
	inputRecords = records
	if err != nil {
		panic(err)
	}

	for _, r := range inputRecords {
		ynabRecord, err := ynabConverter.ConvertFromInternalRecord(r)
		if err != nil {
			panic(err)
		}
		ynabRecords = append(ynabRecords, ynabRecord.(formats.YnabRecord))
	}
	marshalled, err := gocsv.MarshalString(ynabRecords)

	if err != nil {
		panic(err)
	}
	fmt.Print(marshalled)
}
