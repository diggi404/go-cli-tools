package webmail

import (
	"errors"
	"fmt"
	"strings"
)

func ReformCreds(creds string) ([]string, error) {
	splittedCreds := strings.Split(creds, ":")
	if len(splittedCreds) != 2 {
		err := errors.New("invalid credentials format")
		return nil, err
	}
	domain := strings.Split(splittedCreds[0], "@")
	if len(domain) != 2 {
		err := errors.New("invalid credentials format")
		return nil, err
	}
	targetURL := fmt.Sprintf("https://%s:2096", domain[1])
	splittedCreds = append(splittedCreds, targetURL)
	return splittedCreds, nil
}
