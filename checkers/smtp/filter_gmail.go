package smtp

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
		err := errors.New("domain is gmail")
		return nil, err
	}
	return splitedCreds, nil
}
