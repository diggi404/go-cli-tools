package smtpbrute

import (
	"os"
)

func WriteResultsToFile() os.File {
	dirPath, _ := CheckDir()
	filePath := dirPath + "/cracked_smtps.txt"
	file, _ := os.Create(filePath)
	return *file
}
