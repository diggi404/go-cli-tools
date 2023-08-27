package cpanelbrute

import (
	"fmt"
	"os"
)

func SetupDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	dirPath := cwd + "/cpanel_logs"
	_, err = os.Stat(dirPath)
	if err == nil {
		return dirPath, nil
	} else if os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return "", err
		}
		return dirPath, nil
	}
	return "", err
}

func ResultsToFile() os.File {
	dirPath, _ := SetupDir()
	filePath := dirPath + "/valid_logs.txt"
	file, _ := os.Create(filePath)
	return *file
}
