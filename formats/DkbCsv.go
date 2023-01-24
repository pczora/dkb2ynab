package formats

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/charmap"
)

func init() {
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		// DKB uses ISO8859-15 (for whatever reason)
		reader := csv.NewReader(charmap.ISO8859_15.NewDecoder().Reader(in))
		reader.Comma = ';'
		return reader
	})
}

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
	normalizedAmount := amount.normalizeAmount(csv)
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

func (amount *DkbAmount) normalizeAmount(a string) string {
	result := strings.Replace(a, ".", "", -1)
	result = strings.Replace(result, ",", ".", -1)
	return result
}

type DkbFormatConverter struct{}

func (d *DkbFormatConverter) Identify(path string) bool {
	//TODO: implement
	return false
}

func (d *DkbFormatConverter) ConvertFromInternalRecord(i InternalRecord) (Record, error) {

	return DkbRecord{Date: DkbDateTime(i.Date), ValueDate: DkbDateTime(i.ValueDate), PostingText: i.PostingText, Payee: i.Payee, Purpose: i.Purpose, BankAccountNumber: i.BankAccountNumber, BankCode: i.BankCode, Amount: DkbAmount(i.Amount), CreditorID: i.CreditorID, MandateReference: i.MandateReference, CustomerReference: i.CustomerReference}, nil
}

func (d *DkbFormatConverter) ConvertToInternalRecord(r Record) (InternalRecord, error) {
	record, ok := r.(DkbRecord)
	if !ok {
		return InternalRecord{}, errors.New("Record is not of type DkbCreditCardRecord")
	}
	internalRecord := InternalRecord{Date: DateTime(record.Date), ValueDate: DateTime(record.ValueDate), PostingText: record.PostingText, Payee: record.Payee, Purpose: record.Purpose, BankAccountNumber: record.BankAccountNumber, BankCode: record.BankCode, Amount: Amount(record.Amount), CreditorID: record.CreditorID, MandateReference: record.MandateReference, CustomerReference: record.CustomerReference}
	return internalRecord, nil
}

func (d *DkbFormatConverter) ConvertFromFile(path string) ([]InternalRecord, error) {
	f, err := os.Open(path)

	if err != nil {
		return []InternalRecord{}, err
	}

	defer f.Close()

	fileReader := bufio.NewReader(f)

	skipLines(fileReader, 6)

	dkbRecords := []DkbRecord{}
	err = gocsv.Unmarshal(fileReader, &dkbRecords)
	if err != nil {
		return []InternalRecord{}, err
	}
	var result []InternalRecord
	for _, r := range dkbRecords {
		genericRecord, err := d.ConvertToInternalRecord(r)
		if err != nil {
			return result, err
		}
		result = append(result, genericRecord)
	}
	return result, nil
}

type DkbCreditCardRecord struct {
	Marked    string      `csv:"Umsatz abgerechnet aber nicht im Saldo enthalten"` // Ignored (for now)
	ValueDate DkbDateTime `csv:"Wertstellung"`
	Date      DkbDateTime `csv:"Belegdatum"`
	Purpose   string      `csv:"Beschreibung"`
	Amount    DkbAmount   `csv:"Betrag (EUR)"`
	//OriginalAmount DkbAmount   `csv:"Ursprünglicher Betrag"` // Ignored (for now)
}

type DkbCreditCardFormatConverter struct{}

func (d *DkbCreditCardFormatConverter) Identify(path string) bool {
	//TODO: implement
	return false
}

func (d *DkbCreditCardFormatConverter) ConvertFromInternalRecord(i InternalRecord) (Record, error) {
	return DkbCreditCardRecord{"", DkbDateTime(i.ValueDate), DkbDateTime(i.Date), i.Purpose, DkbAmount(i.Amount)}, nil
}

func (d *DkbCreditCardFormatConverter) ConvertToInternalRecord(r Record) (InternalRecord, error) {
	record, ok := r.(DkbCreditCardRecord)
	if !ok {
		return InternalRecord{}, errors.New("Record is not of type DkbCreditCardRecord")
	}
	return InternalRecord{Date: DateTime(record.Date), ValueDate: DateTime(record.ValueDate), PostingText: record.Purpose, Payee: record.Purpose, Purpose: record.Purpose, BankAccountNumber: "", BankCode: "", Amount: Amount(record.Amount), CreditorID: "", MandateReference: "", CustomerReference: ""}, nil
}

func (d *DkbCreditCardFormatConverter) ConvertFromFile(path string) ([]InternalRecord, error) {
	f, err := os.Open(path)

	if err != nil {
		return []InternalRecord{}, errors.New("cannot open file")
	}

	defer f.Close()

	fileReader := bufio.NewReader(f)

	skipLines(fileReader, 6)

	dkbCreditCardRecords := []DkbCreditCardRecord{}
	err = gocsv.Unmarshal(fileReader, &dkbCreditCardRecords)
	if err != nil {
		panic(err)
	}
	var result []InternalRecord
	for _, r := range dkbCreditCardRecords {
		genericRecord, err := d.ConvertToInternalRecord(r)
		if err != nil {
			return result, errors.New("could not convert record")
		}
		result = append(result, genericRecord)
	}
	return result, nil
}

func skipLines(r *bufio.Reader, n int) {
	for i := 0; i < n; i++ {
		r.ReadLine()
	}
}
