package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/pczora/dkb2ynab/formats"
	"github.com/pczora/dkbrobot/pkg/dkbclient"
	"github.com/pczora/zprobot/pkg/zinspilotclient"
	"golang.org/x/term"
)

const (
	DateLayout = "2006-01-02"
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

	dkbAccounts, err := dkb.GetAccounts()
	if err != nil {
		panic(err)
	}

	for _, a := range dkbAccounts.Data {
		createDkbCheckingAccountCsv(&dkb, a)
	}

	dkbCreditCards, err := dkb.GetCreditCards()
	if err != nil {
		panic(err)
	}

	for _, c := range dkbCreditCards.Data {
		if c.Type == "creditCard" {
			createDkbCreditCardCsv(&dkb, c)
		}
	}
}

func createDkbCheckingAccountCsv(dkb *dkbclient.Client, account dkbclient.Account) {

	var ynabConverter formats.YnabFormatConverter
	var records []formats.InternalRecord
	var ynabRecords []formats.YnabRecord

	transactions, err := dkb.GetAccountTransactions(account.Id)
	if err != nil {
		fmt.Println(err)
	}

	for _, t := range transactions.Data {
		date, err := time.Parse(DateLayout, t.Attributes.BookingDate)
		if err != nil {
			panic(err)
		}

		if t.Attributes.Status == "pending" || t.Attributes.ValueDate == "" {
			fmt.Println("Transaction pending, skipping...")
			continue
		}

		valueDate, err := time.Parse(DateLayout, t.Attributes.ValueDate)
		if err != nil {
			panic(err)
		}

		amount, err := strconv.ParseFloat(t.Attributes.Amount.Value, 64)
		if err != nil {
			panic(err)
		}
		r := formats.InternalRecord{}

		if amount < 0 {
			r = formats.InternalRecord{Date: date, ValueDate: valueDate, PostingText: t.Attributes.Description, Payee: t.Attributes.Creditor.Name, Purpose: t.Attributes.Description, BankAccountNumber: t.Attributes.Creditor.CreditorAccount.Iban, BankCode: t.Attributes.Creditor.CreditorAccount.Blz, Amount: amount, CreditorID: t.Attributes.Creditor.Id, MandateReference: t.Attributes.MandateId, CustomerReference: t.Attributes.EndToEndId}
		} else {
			r = formats.InternalRecord{Date: date, ValueDate: valueDate, PostingText: t.Attributes.Description, Payee: t.Attributes.Debtor.Name, Purpose: t.Attributes.Description, BankAccountNumber: t.Attributes.Debtor.DebtorAccount.Iban, BankCode: t.Attributes.Debtor.DebtorAccount.Blz, Amount: amount, CreditorID: t.Attributes.Creditor.Id, MandateReference: t.Attributes.MandateId, CustomerReference: t.Attributes.EndToEndId}

		}
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

	err = os.MkdirAll(fmt.Sprintf("output/dkb/%v/%v/", account.Attributes.HolderName, account.Type), os.ModePerm)
	if err != nil {
		return
	}
	err = os.WriteFile(fmt.Sprintf("output/dkb/%v/%v/%v.csv", account.Attributes.HolderName, account.Type, account.Attributes.Iban), marshalled, 0644)
	if err != nil {
		panic(err)
	}
}

func createDkbCreditCardCsv(dkb *dkbclient.Client, c dkbclient.CreditCard) {

	var ynabConverter formats.YnabFormatConverter
	var records []formats.InternalRecord
	var ynabRecords []formats.YnabRecord

	transactions, err := dkb.GetCreditCardTransactions(c.Id)
	if err != nil {
		fmt.Println(err)
	}

	for _, t := range transactions.Data {
		bookingDate, err := time.Parse(DateLayout, t.Attributes.BookingDate)
		if err != nil {
			panic(err)
		}

		amount, err := strconv.ParseFloat(t.Attributes.Amount.Value, 64)
		if err != nil {
			panic(err)
		}

		r := formats.InternalRecord{Date: t.Attributes.AuthorizationDate, ValueDate: bookingDate, PostingText: t.Attributes.Description, Payee: t.Attributes.Description, Purpose: t.Attributes.Description, BankAccountNumber: "", BankCode: "", Amount: amount, CreditorID: "", MandateReference: "", CustomerReference: ""}
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
	err = os.MkdirAll(fmt.Sprintf("output/dkb/%v %v/%v/", c.Attributes.Owner.FirstName, c.Attributes.Owner.LastName, c.Type), os.ModePerm)
	if err != nil {
		return
	}
	err = os.WriteFile(fmt.Sprintf("output/dkb/%v %v/%v/%v.csv", c.Attributes.Owner.FirstName, c.Attributes.Owner.LastName, c.Type, c.Attributes.MaskedPan), marshalled, 0644)
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
