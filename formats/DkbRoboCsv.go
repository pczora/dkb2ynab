package formats

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
)

type DkbRoboConverter struct{}

func init() {
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		return gocsv.LazyCSVReader(in) // Allows use of quotes in CSV
	})

}

type DkbRoboDateTime struct {
	time.Time
}

func (date *DkbRoboDateTime) MarshalCSV() (string, error) {
	return date.Time.Format("02.01.2006"), nil
}

func (date *DkbRoboDateTime) UnmarshalCSV(csv string) (err error) {
	t, err := time.Parse("02.01.2006", csv)
	date.Time = t
	return err
}

type DkbRoboAmount struct {
	float64
}

func (amount *DkbRoboAmount) MarshalCSV() (string, error) {
	return strconv.FormatFloat(amount.float64, 'f', 2, 64), nil
}

func (amount *DkbRoboAmount) UnmarshalCSV(csv string) (err error) {
	if csv == "" {
		amount.float64 = 0.0
		return nil
	}
	floatAmount, err := strconv.ParseFloat(csv, 64)
	if err != nil {
		return err
	}
	amount.float64 = floatAmount
	return nil
}

type DkbRoboRecord struct {
	BDate             DkbRoboDateTime `csv:"bdate"`
	VDate             DkbRoboDateTime `csv:"vdate"`
	PostingText       string          `csv:"postingtext"`
	Peer              string          `csv:"peer"`
	ReasonForPayment  string          `csv:"reasonforpayment"`
	MandateReference  string          `csv:"mandatereference"`
	CustomerReference string          `csv:"customerreferenz"`
	PeerAccount       string          `csv:"peeraccount"`
	PeerBIC           string          `csv:"peerbic"`
	PeerID            string          `csv:"peerid"`
	Amount            DkbRoboAmount   `csv:"amount"`
	Date              DkbRoboDateTime `csv:"date"`
	Text              string          `csv:"text"`
}

func (d *DkbRoboConverter) ConvertFromInternalRecord(i InternalRecord) (Record, error) {
	panic("not implemented") // TODO: Implement
}

func (d *DkbRoboConverter) ConvertToInternalRecord(r Record) (InternalRecord, error) {
	record, ok := r.(DkbRoboRecord)
	if !ok {
		return InternalRecord{}, errors.New("Record is not of type DkbRoboRecord")
	}
	internalRecord := InternalRecord{Date: DateTime(record.BDate), ValueDate: DateTime(record.VDate), PostingText: record.PostingText, Payee: record.Peer, Purpose: record.Text, BankAccountNumber: record.PeerAccount, BankCode: record.PeerBIC, Amount: Amount(record.Amount), CreditorID: record.PeerID, MandateReference: record.MandateReference, CustomerReference: record.CustomerReference}
	return internalRecord, nil
}

func (d *DkbRoboConverter) ConvertFromFile(path string) ([]InternalRecord, error) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("Could not open file")
		return []InternalRecord{}, err
	}

	defer f.Close()

	fileReader := bufio.NewReader(f)

	dkbRoboRecords := []DkbRoboRecord{}
	err = gocsv.Unmarshal(fileReader, &dkbRoboRecords)
	if err != nil {
		fmt.Println("Could not unmarshal input")
		return []InternalRecord{}, err
	}
	var result []InternalRecord
	for _, r := range dkbRoboRecords {
		genericRecord, err := d.ConvertToInternalRecord(r)
		if err != nil {
			return result, err
		}
		result = append(result, genericRecord)
	}
	return result, nil
}

func (d *DkbRoboConverter) Identify(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Could not open file")
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	return scanner.Text() == "bdate,vdate,postingtext,peer,reasonforpayment,mandatereference,customerreferenz,peeraccount,peerbic,peerid,amount,date,text"
}

type DkbRoboCreditCardRecord struct {
	BDate          DkbRoboDateTime `csv:"bdate"`
	VDate          DkbRoboDateTime `csv:"vdate"`
	ShowDate       DkbRoboDateTime `csv:"show_date"`
	StoreDate      DkbRoboDateTime `csv:"store_date"`
	Amount         DkbRoboAmount   `csv:"amount"`
	AmountOriginal DkbRoboAmount   `csv:"amount_original"`
	Text           string          `csv:"text"`
}

type DkbRoboCreditCardConverter struct{}

func (d *DkbRoboCreditCardConverter) ConvertFromInternalRecord(i InternalRecord) (Record, error) {
	panic("not implemented") // TODO: Implement
}

func (d *DkbRoboCreditCardConverter) ConvertToInternalRecord(r Record) (InternalRecord, error) {
	record, ok := r.(DkbRoboCreditCardRecord)
	if !ok {
		return InternalRecord{}, errors.New("Record is not of type DkbRoboCreditCardRecord")
	}
	internalRecord := InternalRecord{Date: DateTime(record.BDate), ValueDate: DateTime(record.VDate), PostingText: "", Payee: record.Text, Purpose: record.Text, BankAccountNumber: "", BankCode: "", Amount: Amount(record.Amount), CreditorID: "", MandateReference: "", CustomerReference: ""}
	return internalRecord, nil
}

func (d *DkbRoboCreditCardConverter) ConvertFromFile(path string) ([]InternalRecord, error) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("Could not open file")
		return []InternalRecord{}, err
	}

	defer f.Close()

	fileReader := bufio.NewReader(f)

	dkbRoboCreditCardRecords := []DkbRoboCreditCardRecord{}
	err = gocsv.Unmarshal(fileReader, &dkbRoboCreditCardRecords)
	if err != nil {
		fmt.Println("Could not unmarshal input")
		return []InternalRecord{}, err
	}
	var result []InternalRecord
	for _, r := range dkbRoboCreditCardRecords {
		genericRecord, err := d.ConvertToInternalRecord(r)
		if err != nil {
			return result, err
		}
		result = append(result, genericRecord)
	}
	return result, nil
}

func (d *DkbRoboCreditCardConverter) Identify(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Could not open file")
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	return scanner.Text() == "vdate,show_date,bdate,store_date,text,amount,amount_original"
}
