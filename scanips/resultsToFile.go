package scanips

import (
	"os"
)

func WriteToFile() os.File {
	dirPath, _ := SetupDir()
	filePath := dirPath + "/scanned_ips.txt"
	file, _ := os.Create(filePath)
	return *file
}
