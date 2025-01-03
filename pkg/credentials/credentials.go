package credentials

import (
	"fmt"
	"github.com/zalando/go-keyring"
	"golang.org/x/term"
	"syscall"
)

const CredentialPrefix = "dkb2ynab_"

func FromKeyring(account, username string) (string, error) {
	key := fmt.Sprintf("%v%v", CredentialPrefix, account)
	password, err := keyring.Get(key, username)
	if err != nil {
		return "", err
	}
	return password, nil
}

func FromInteractiveInput(bank, username string) (string, error) {
	fmt.Printf("%v password for user %v: ", bank, username)
	pwBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	fmt.Print("\n")

	return string(pwBytes), nil
}
