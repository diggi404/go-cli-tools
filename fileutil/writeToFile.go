package fileutil

import (
	"fmt"
	"os"
)

func WriteToFile(dirName, fileName string) *os.File {
	dirPath, err := SetupDir(dirName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	filePath := dirPath + "/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return file
}
