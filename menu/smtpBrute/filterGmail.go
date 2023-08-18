package smtpbrute

import (
	"errors"
	"strings"
)

func FilterGmailCreds(creds string) ([]string, error) {
	splitedCreds := strings.Split(creds, ":")
	username, password := splitedCreds[0], splitedCreds[1]
	if strings.Contains(username, "@gmail.com") {
		err := errors.New("domain is gmail")
		splitedCreds = []string{username, password, "smtp.gmail.com"}
		return splitedCreds, err
	}
	return splitedCreds, nil
}
