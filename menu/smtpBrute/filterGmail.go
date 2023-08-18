package smtpbrute

import (
	"errors"
	"strings"
)

func FilterGmailCreds(creds string) ([]string, error) {
	splitedCreds := strings.Split(creds, ":")
	if len(splitedCreds) != 2 {
		err := errors.New("please make sure your crendentials are separated by colon (:)")
		return nil, err
	}
	username, password := splitedCreds[0], splitedCreds[1]
	if strings.Contains(username, "@gmail.com") {
		err := errors.New("domain is gmail")
		splitedCreds = []string{username, password, "smtp.gmail.com"}
		return splitedCreds, err
	}
	return splitedCreds, nil
}
