package formats

import (
	"strconv"
	"time"
)

type YnabDateTime struct {
	DateTime
}

func (date *YnabDateTime) MarshalCSV() (string, error) {
	return date.Time.Format("2006/01/02"), nil
}

func (date *YnabDateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2006/01/02", csv)
	return err
}

type YnabAmount struct {
	Amount
}

func (amount *YnabAmount) MarshalCSV() (string, error) {
	return strconv.FormatFloat(amount.amount, 'f', 2, 64), nil
}

func (amount *YnabAmount) UnmarshalCSV(csv string) (err error) {
	floatAmount, err := strconv.ParseFloat(csv, 64)
	if err != nil {
		return err
	}
	amount.Amount.amount = floatAmount
	return nil
}

type YnabRecord struct {
	Date   YnabDateTime `csv:"Date"`
	Payee  string       `csv:"Payee"`
	Memo   string       `csv:"Memo"`
	Amount YnabAmount   `csv:"Amount"`
}
