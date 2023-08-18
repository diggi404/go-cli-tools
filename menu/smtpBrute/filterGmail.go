package smtpbrute

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func FilterGmailCreds(creds string) ([]string, error) {
	splitedCreds := strings.Split(creds, ":")
	if len(splitedCreds) != 2 {
		fmt.Println("Please make sure your wordlist is separated by a colon (:)")
		os.Exit(2)
	}
	username, password := splitedCreds[0], splitedCreds[1]
	if strings.Contains(username, "@gmail.com") {
		err := errors.New("domain is gmail")
		splitedCreds = []string{username, password, "smtp.gmail.com"}
		return splitedCreds, err
	}
	return splitedCreds, nil
}
