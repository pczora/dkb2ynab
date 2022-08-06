package formats

import "time"

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
	ConvertFromFile(path string) ([]InternalRecord, error)
}
