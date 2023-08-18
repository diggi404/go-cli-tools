package smtpbrute

import (
	"fmt"
	"os"
)

func WriteResultsToFile(results []string) {
	dirPath, err := CheckDir()
	if err != nil {
		return
	}
	filePath := dirPath + "/cracked_smtps.txt"
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	for _, result := range results {
		_, err = file.WriteString(result)
		if err != nil {
			return
		}
	}
}
