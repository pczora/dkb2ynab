package formats

import (
	"errors"
	"time"
)

var converters = []Converter{&DkbRoboConverter{}, &YnabFormatConverter{}, &DkbCreditCardFormatConverter{}, &DkbFormatConverter{}, &DkbRoboCreditCardConverter{}, &ZinspilotFormatConverter{}}

type DateTime struct {
	time.Time
}

type Amount struct {
	float64
}

type InternalRecord struct {
	Date              DateTime `csv:"Date"`
	ValueDate         DateTime `csv:"Value Date"`
	PostingText       string   `csv:"Posting Text"`
	Payee             string   `csv:"Payee"`
	Purpose           string   `csv:"Purpose"`
	BankAccountNumber string   `csv:"Bank Account Number"`
	BankCode          string   `csv:"Bank Code"`
	Amount            Amount   `csv:"Amount"`
	CreditorID        string   `csv:"Creditor ID"`
	MandateReference  string   `csv:"Mandate Reference"`
	CustomerReference string   `csv:"Customer Reference"`
}

type Record interface{}

type Converter interface {
	ConvertFromInternalRecord(i InternalRecord) (Record, error)
	ConvertToInternalRecord(r Record) (InternalRecord, error)
	//TODO: This should rather convert from a byte array/a reader/... than from a file
	ConvertFromFile(path string) ([]InternalRecord, error)
	Identify(path string) bool
}

func FindSuitableConverter(path string) (Converter, error) {
	for _, converter := range converters {
		if converter.Identify(path) {
			return converter, nil
		}
	}
	return nil, errors.New("could not find suitable converter")
}
