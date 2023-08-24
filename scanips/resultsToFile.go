package scanips

import (
	"fmt"
	"os"
)

func WriteToFile(results []string) {
	dirPath, err := SetupDir()
	if err != nil {
		return
	}
	filePath := dirPath + "/scanned_ips.txt"
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
