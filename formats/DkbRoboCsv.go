package formats

import (
	"bufio"
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

func (d *DkbRoboConverter) ConvertToInternalRecord(r DkbRoboRecord) (InternalRecord, error) {
	internalRecord := InternalRecord{Date: DateTime(r.BDate), ValueDate: DateTime(r.VDate), PostingText: r.PostingText, Payee: r.Peer, Purpose: r.Text, BankAccountNumber: r.PeerAccount, BankCode: r.PeerBIC, Amount: Amount(r.Amount), CreditorID: r.PeerID, MandateReference: r.MandateReference, CustomerReference: r.CustomerReference}
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
