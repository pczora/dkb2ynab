package formats

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

func init() {
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		return gocsv.LazyCSVReader(in) // Allows use of quotes in CSV
	})
}

type ZinspilotDayTime struct {
	time.Time
}

func (date *ZinspilotDayTime) MarshalCSV() (string, error) {
	return date.Format("02.01.2006"), nil
}

func (date *ZinspilotDayTime) UnmarshalCSV(csv string) (err error) {
	t, err := time.Parse("02.01.2006", csv)
	date.Time = t
	return err
}

type ZinspilotAmount struct {
	float64
}

func (amount *ZinspilotAmount) MarshalCSV() (string, error) {
	return strconv.FormatFloat(amount.float64, 'f', 2, 64), nil
}

func (amount *ZinspilotAmount) UnmarshalCSV(csv string) (err error) {
	normalizedAmount := amount.normalizeAmount(csv)
	floatAmount, err := strconv.ParseFloat(normalizedAmount, 64)
	if err != nil {
		return err
	}
	amount.float64 = floatAmount
	return nil
}

func (amount *ZinspilotAmount) normalizeAmount(a string) string {
	result := strings.ReplaceAll(a, ".", "")
	result = strings.ReplaceAll(result, ",", ".")
	result = strings.ReplaceAll(result, "\"", "")
	return result
}

type ZinspilotRecord struct {
	Buchungsdatum ZinspilotDayTime `csv:"Buchungsdatum"`
	Wertstellung  ZinspilotDayTime `csv:"Wertstellung"`
	Umsatz        string           `csv:"Umsatz"`
	Saldo         ZinspilotAmount  `csv:"Saldo"`
}

type ZinspilotFormatConverter struct {
}

func (z *ZinspilotFormatConverter) Identify(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Could not open file")
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	return scanner.Text() == "Buchungsdatum,Wertstellung,Umsatz,Saldo"
}

func (z ZinspilotFormatConverter) ConvertFromFile(path string) ([]InternalRecord, error) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("Could not open file")
		return []InternalRecord{}, err
	}

	defer f.Close()

	fileReader := bufio.NewReader(f)

	zinspilotRecords := []ZinspilotRecord{}
	err = gocsv.Unmarshal(fileReader, &zinspilotRecords)
	if err != nil {
		fmt.Println("Could not unmarshal input")
		return []InternalRecord{}, err
	}
	var result []InternalRecord
	for _, r := range zinspilotRecords {
		genericRecord, err := z.ConvertToInternalRecord(r)
		if err != nil {
			return result, err
		}
		result = append(result, genericRecord)
	}
	return result, nil
}

func (z ZinspilotFormatConverter) ConvertFromInternalRecord(r InternalRecord) (Record, error) {
	zinspilotRecord := ZinspilotRecord{Buchungsdatum: ZinspilotDayTime(r.Date), Umsatz: r.Purpose, Saldo: ZinspilotAmount(r.Amount)}
	return zinspilotRecord, nil
}

func (z ZinspilotFormatConverter) ConvertToInternalRecord(r Record) (InternalRecord, error) {
	record, ok := r.(ZinspilotRecord)
	if !ok {
		return InternalRecord{}, errors.New("Record is not of type ZinspilotRecord")
	}
	internalRecord := InternalRecord{Date: DateTime(record.Buchungsdatum), ValueDate: DateTime(record.Wertstellung), Payee: "", Purpose: record.Umsatz, PostingText: record.Umsatz, Amount: Amount(record.Saldo)}
	return internalRecord, nil
}
