package scanips

import (
	"bufio"
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()
		ips = append(ips, ip)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return ips, nil
}
