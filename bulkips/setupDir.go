package bulkips

import (
	"fmt"
	"os"
)

// This function checks if the default directory of the application is available and if not it makes a new dir.
func CheckDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	dirPath := cwd + "/ip_list"
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
