package smtpbrute

import (
	"fmt"
	"os"
)

func ReadCredsFromFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wordList []string
	var creds string
	for {
		_, err := fmt.Fscanf(file, "%s", &creds)
		if err != nil {
			break
		}
		wordList = append(wordList, creds)
	}

	return wordList, nil
}
