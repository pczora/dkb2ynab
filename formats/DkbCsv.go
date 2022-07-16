package formats

import (
	"strconv"
	"strings"
	"time"
)

type DkbDateTime struct {
	DateTime
}

func (date *DkbDateTime) MarshalCSV() (string, error) {
	return date.Time.Format("02.01.2006"), nil
}

func (date *DkbDateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("02.01.2006", csv)
	return err
}

type DkbAmount struct {
	Amount
}

func (amount *DkbAmount) MarshalCSV() (string, error) {
	return strconv.FormatFloat(amount.amount, 'f', 2, 64), nil
}

func (amount *DkbAmount) UnmarshalCSV(csv string) (err error) {
	normalizedAmount := normalizeAmount(csv)
	floatAmount, err := strconv.ParseFloat(normalizedAmount, 64)
	if err != nil {
		return err
	}
	amount.Amount.amount = floatAmount
	return nil
}

type DkbRecord struct {
	Date              DkbDateTime `csv:"Buchungstag"`
	ValueDate         DkbDateTime `csv:"Wertstellungstag"`
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
