package cpanel

import (
	"errors"
	"strings"
)

func FilterCreds(wordlist string) ([]string, error) {
	splittedCreds := strings.Split(wordlist, "|")
	if len(splittedCreds) != 3 {
		err := errors.New("invalid wordlist format")
		return nil, err
	}
	return splittedCreds, nil
}
