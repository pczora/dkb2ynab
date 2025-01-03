package formats

import (
	"errors"
	"strconv"
	"time"
)

type YnabDateTime time.Time

func (date *YnabDateTime) MarshalCSV() (string, error) {
	return time.Time(*date).Format("2006/01/02"), nil
}

func (date *YnabDateTime) UnmarshalCSV(csv string) (err error) {
	t, err := time.Parse("2006/01/02", csv)
	*date = YnabDateTime(t)
	return err
}

type YnabAmount float64

func (amount *YnabAmount) MarshalCSV() (string, error) {
	return strconv.FormatFloat(float64(*amount), 'f', 2, 64), nil
}

func (amount *YnabAmount) UnmarshalCSV(csv string) (err error) {
	floatAmount, err := strconv.ParseFloat(csv, 64)
	if err != nil {
		return err
	}
	*amount = YnabAmount(floatAmount)
	return nil
}

type YnabRecord struct {
	Date   YnabDateTime `csv:"Date"`
	Payee  string       `csv:"Payee"`
	Memo   string       `csv:"Memo"`
	Amount YnabAmount   `csv:"Amount"`
}

type YnabFormatConverter struct{}

func (y *YnabFormatConverter) Identify(path string) bool {
	//TODO: implement
	return false
}

func (y YnabFormatConverter) ConvertFromFile(path string) ([]InternalRecord, error) {
	panic("not implemented") // TODO: Implement
}

func (y YnabFormatConverter) ConvertFromInternalRecord(r InternalRecord) (Record, error) {
	ynabRecord := YnabRecord{Date: YnabDateTime(r.Date), Payee: r.Payee, Memo: r.Purpose, Amount: YnabAmount(r.Amount)}
	return ynabRecord, nil
}

func (y YnabFormatConverter) ConvertToInternalRecord(r Record) (InternalRecord, error) {
	record, ok := r.(YnabRecord)
	if !ok {
		return InternalRecord{}, errors.New("Record is not of type YnabRecord")
	}
	internalRecord := InternalRecord{Date: time.Time(record.Date), ValueDate: time.Time(record.Date), Payee: record.Payee, PostingText: record.Memo, Amount: float64(record.Amount)}
	return internalRecord, nil
}
