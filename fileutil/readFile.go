package fileutil

import (
	"bufio"
	"os"
)

func ReadFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fileContent []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		creds := scanner.Text()
		fileContent = append(fileContent, creds)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return fileContent, nil
}
