package smtpbrute

import (
	"bufio"
	"os"
)

func ReadCredsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wordList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		creds := scanner.Text()
		wordList = append(wordList, creds)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return wordList, nil
}
