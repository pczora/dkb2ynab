package formats

import (
	"time"
)

type DateTime time.Time

type Amount float64

type InternalRecord struct {
	Date              time.Time `csv:"Date"`
	ValueDate         time.Time `csv:"Value Date"`
	PostingText       string    `csv:"Posting Text"`
	Payee             string    `csv:"Payee"`
	Purpose           string    `csv:"Purpose"`
	BankAccountNumber string    `csv:"Bank Account Number"`
	BankCode          string    `csv:"Bank Code"`
	Amount            float64   `csv:"Amount"`
	CreditorID        string    `csv:"Creditor ID"`
	MandateReference  string    `csv:"Mandate Reference"`
	CustomerReference string    `csv:"Customer Reference"`
}

type Record interface{}

type Converter interface {
	ConvertFromInternalRecord(i InternalRecord) (Record, error)
	ConvertToInternalRecord(r Record) (InternalRecord, error)
	//TODO: This should rather convert from a byte array/a reader/... than from a file
	ConvertFromFile(path string) ([]InternalRecord, error)
	Identify(path string) bool
}
