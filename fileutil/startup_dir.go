package fileutil

import (
	"os"
	"path/filepath"
)

func SetupDir(dirName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dirPath := filepath.Join(cwd, dirName)
	_, err = os.Stat(dirPath)
	if err == nil {
		return dirPath, nil
	} else if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return "", err
		}
		return dirPath, nil
	}
	return "", err
}
