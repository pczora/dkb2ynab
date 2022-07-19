package formats

import (
	"strconv"
	"strings"
	"time"
)

type DkbDateTime struct {
	time.Time
}

func (date *DkbDateTime) MarshalCSV() (string, error) {
	return date.Time.Format("02.01.2006"), nil
}

func (date *DkbDateTime) UnmarshalCSV(csv string) (err error) {
	t, err := time.Parse("02.01.2006", csv)
	date.Time = t
	return err
}

type DkbAmount struct {
	float64
}

func (amount *DkbAmount) MarshalCSV() (string, error) {
	return strconv.FormatFloat(amount.float64, 'f', 2, 64), nil
}

func (amount *DkbAmount) UnmarshalCSV(csv string) (err error) {
	normalizedAmount := normalizeAmount(csv)
	floatAmount, err := strconv.ParseFloat(normalizedAmount, 64)
	if err != nil {
		return err
	}
	amount.float64 = floatAmount
	return nil
}

type DkbRecord struct {
	Date              DkbDateTime `csv:"Buchungstag"`
	ValueDate         DkbDateTime `csv:"Wertstellung"`
	PostingText       string      `csv:"Buchungstext"`
	Payee             string      `csv:"Auftraggeber / Begünstigter"`
	Purpose           string      `csv:"Verwendungszweck"`
	BankAccountNumber string      `csv:"Kontonummer"`
	BankCode          string      `csv:"Bankleitzahl"`
	Amount            DkbAmount   `csv:"Betrag (EUR)"`
	CreditorID        string      `csv:"Gläubiger-ID"`
	MandateReference  string      `csv:"Mandatsreferenz"`
	CustomerReference string      `csv:"Kundenreferenz"`
}

func normalizeAmount(amount string) string {
	result := strings.Replace(amount, ".", "", -1)
	result = strings.Replace(result, ",", ".", -1)
	return result
}

type DkbFormatConverter struct{}

func (d *DkbFormatConverter) ConvertFromInternalRecord(r InternalRecord) DkbRecord {
	return DkbRecord{Date: DkbDateTime(r.Date), ValueDate: DkbDateTime(r.ValueDate), PostingText: r.PostingText, Payee: r.Payee, Purpose: r.Purpose, BankAccountNumber: r.BankAccountNumber, BankCode: r.BankCode, Amount: DkbAmount(r.Amount), CreditorID: r.CreditorID, MandateReference: r.MandateReference, CustomerReference: r.CustomerReference}
}

func (d *DkbFormatConverter) ConvertToInternalRecord(r DkbRecord) InternalRecord {
	internalRecord := InternalRecord{Date: DateTime(r.Date), ValueDate: DateTime(r.ValueDate), PostingText: r.PostingText, Payee: r.Payee, Purpose: r.Purpose, BankAccountNumber: r.BankAccountNumber, BankCode: r.BankCode, Amount: Amount(r.Amount), CreditorID: r.CreditorID, MandateReference: r.MandateReference, CustomerReference: r.CustomerReference}
	return internalRecord
}
