package menu

import (
	"fmt"
	"os"
)

// ReadIPsFromFile This reads all the IPs from the path passed to it.
func ReadIPsFromFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ips []string
	var ip string
	for {
		_, err := fmt.Fscanf(file, "%s", &ip)
		if err != nil {
			break
		}
		ips = append(ips, ip)
	}

	return ips, nil
}
