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
	//var dkbCreditCardRecords []formats.InternalRecord
	//dkbConverter := formats.DkbFormatConverter{}
	dkbRoboConverter := formats.DkbRoboConverter{}
	//dkbCreditCardConverter := formats.DkbCreditCardFormatConverter{}
	ynabConverter := formats.YnabFormatConverter{}
	//dkbRecords, err := dkbConverter.ConvertFromFile(args[1])
	dkbRoboRecords, err := dkbRoboConverter.ConvertFromFile(args[1])
	if err != nil {
		panic(err)
	}
	for _, r := range dkbRoboRecords {
		ynabRecord, err := ynabConverter.ConvertFromInternalRecord(r)
		if err != nil {
			panic(err)
		}
		ynabRecords = append(ynabRecords, ynabRecord)
	}
	//dkbCreditCardRecords, err := dkbCreditCardConverter.ConvertFromFile(args[1])
	//for _, r := range dkbCreditCardRecords {
	//ynabRecord, err := ynabConverter.ConvertFromInternalRecord(r)
	//if err != nil {
	//panic("Could not convert record")
	//}
	//ynabRecords = append(ynabRecords, ynabRecord)
	//}
	marshalled, err := gocsv.MarshalString(ynabRecords)

	if err != nil {
		panic(err)
	}
	fmt.Print(marshalled)
}
