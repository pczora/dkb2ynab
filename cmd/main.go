package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/pczora/dkb2ynab/formats"
	"github.com/pczora/dkbrobot/pkg/dkbclient"
	"github.com/pczora/zprobot/pkg/zinspilotclient"
	"golang.org/x/term"
)

func main() {
	fetchDkbTransactions()
// 	fetchZinspilotTransactions()
}

func fetchDkbTransactions() {
	var dkbUsername string
	var dkbPassword string

	fmt.Printf("DKB username: ")
	_, err := fmt.Scanf("%s", &dkbUsername)
	if err != nil {
		panic(err)
	}

	fmt.Printf("DKB password: ")
	bytepw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		os.Exit(1)
	}
	fmt.Print("\n")

	dkbPassword = string(bytepw)

	dkb := dkbclient.New()

	err = dkb.Login(dkbUsername, dkbPassword)
	if err != nil {
		panic(err)
	}

	dkbAccounts, err := dkb.ParseOverview()
	if err != nil {
		panic(err)
	}

	for _, a := range dkbAccounts {
		switch a.AccountType {
		case dkbclient.CheckingAccount:
			createDkbCheckingAccountCsv(&dkb, a)
		case dkbclient.CreditCard:
			createDkbCreditCardCsv(&dkb, a)
		}
		if a.AccountType != dkbclient.CheckingAccount {
			continue
		}

	}
}

func createDkbCheckingAccountCsv(dkb *dkbclient.Client, a dkbclient.AccountMetadata) {

	var ynabConverter formats.YnabFormatConverter
	var records []formats.InternalRecord
	var ynabRecords []formats.YnabRecord

	transactions, err := dkb.GetAccountTransactions(a, time.Now().Add(30*-time.Hour*24), time.Now())
	if err != nil {
		fmt.Println(err)
	}

	for _, t := range transactions {
		r := formats.InternalRecord{Date: time.Time(t.Date), ValueDate: time.Time(t.ValueDate), PostingText: t.PostingText, Payee: t.Payee, Purpose: t.Purpose, BankAccountNumber: t.BankAccountNumber, BankCode: t.BankCode, Amount: float64(t.Amount), CreditorID: t.CreditorID, MandateReference: t.MandateReference, CustomerReference: t.CustomerReference}
		records = append(records, r)
	}

	for _, r := range records {
		y, err := ynabConverter.ConvertFromInternalRecord(r)
		if err != nil {
			fmt.Println(err)
		}
		ynabRecords = append(ynabRecords, y.(formats.YnabRecord))
	}

	marshalled, err := gocsv.MarshalBytes(ynabRecords)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(fmt.Sprintf("output/dkb/%v.csv", a.Account), marshalled, 0644)
	if err != nil {
		panic(err)
	}
}

func createDkbCreditCardCsv(dkb *dkbclient.Client, a dkbclient.AccountMetadata) {

	var ynabConverter formats.YnabFormatConverter
	var records []formats.InternalRecord
	var ynabRecords []formats.YnabRecord

	transactions, err := dkb.GetCreditCardTransactions(a, time.Now().Add(30*-time.Hour*24), time.Now())
	if err != nil {
		fmt.Println(err)
	}

	for _, t := range transactions {
		r := formats.InternalRecord{Date: time.Time(t.Date), ValueDate: time.Time(t.ValueDate), PostingText: t.Purpose, Payee: t.Purpose, Purpose: t.Purpose, BankAccountNumber: "", BankCode: "", Amount: float64(t.Amount), CreditorID: "", MandateReference: "", CustomerReference: ""}
		records = append(records, r)
	}

	for _, r := range records {
		y, err := ynabConverter.ConvertFromInternalRecord(r)
		if err != nil {
			fmt.Println(err)
		}
		ynabRecords = append(ynabRecords, y.(formats.YnabRecord))
	}

	marshalled, err := gocsv.MarshalBytes(ynabRecords)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(fmt.Sprintf("output/dkb/%v.csv", a.Account), marshalled, 0644)
	if err != nil {
		panic(err)
	}
}

func fetchZinspilotTransactions() {
	var username string
	var password string
	var records []formats.InternalRecord
	var ynabRecords []formats.YnabRecord
	var ynabConverter formats.YnabFormatConverter

	fmt.Printf("Zinspilot username: ")
	_, err := fmt.Scanf("%s", &username)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Zinspilot password: ")
	bytepw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		os.Exit(1)
	}
	fmt.Print("\n")

	password = string(bytepw)

	zp := zinspilotclient.New()

	err = zp.Login(username, password)
	if err != nil {
		panic(err)
	}

	accounts, err := zp.GetAccounts()
	if err != nil {
		panic(err)
	}

	for _, a := range accounts {
		transactions, err := zp.GetTransactions(a)
		if err != nil {
			panic(err)
		}

		for _, t := range transactions {
			if t.Saldo != 0 {
				r := formats.InternalRecord{Date: t.Buchungsdatum, ValueDate: t.Wertstellung, Payee: "Zinsbuchung", Purpose: t.Umsatz, PostingText: t.Umsatz, Amount: t.Saldo}
				records = append(records, r)
			}
		}
		for _, r := range records {
			y, err := ynabConverter.ConvertFromInternalRecord(r)
			if err != nil {
				fmt.Println(err)
			}
			ynabRecords = append(ynabRecords, y.(formats.YnabRecord))
		}

		marshalled, err := gocsv.MarshalBytes(ynabRecords)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(fmt.Sprintf("output/zinspilot/%v.csv", a.Name), marshalled, 0644)
		if err != nil {
			panic(err)
		}

	}

}
