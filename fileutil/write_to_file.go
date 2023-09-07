package fileutil

import (
	"os"
)

func WriteToFile(dirName, fileName string) (*os.File, error) {
	dirPath, err := SetupDir(dirName)
	if err != nil {
		return nil, err

	}
	filePath := dirPath + "/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
