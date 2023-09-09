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
	username := splitedCreds[0]
	if strings.Contains(username, "@gmail.com") {
		splitedCreds = append(splitedCreds, "smtp.gmail.com")
		return splitedCreds, nil
	}
	return splitedCreds, nil
}
