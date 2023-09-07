package bulkips

import (
	"errors"
	"fmt"
	"go_cli/bomber"
	"net"
	"os"
	"time"
)

// GenIP main function to generate the bulk IPs.
func GenIP() (string, error) {
	var ip1, ip2 string
	fmt.Print("\n\nEnter the starting IP: ")
	fmt.Scanln(&ip1)
	startIP := net.ParseIP(ip1).To4()
	if startIP == nil {
		err := errors.New("invalid IP address")
		return "", err
	}
	fmt.Print("\nEnter the ending IP: ")
	fmt.Scanln(&ip2)
	endIP := net.ParseIP(ip2).To4()
	if endIP == nil {
		err := errors.New("invalid IP address")
		return "", err
	}
	fmt.Println()

	// initialize the bar cli output
	totalIps := CountTotalIPs(startIP, endIP)
	pgBar := bomber.MakePgBar(totalIps, "Generating...")

	// prepare the result file's directory and file.
	dirPath, err := CheckDir()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "", err
	}
	currentTime := time.Now().Unix()
	fileName := fmt.Sprintf("generated_ips_%v.txt", currentTime)
	filePath := dirPath + "/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return "", err
	}
	defer file.Close()

	// This loop generates IPs with the range specified by the user.
	for ip := startIP; ip != nil && BytesCompare(ip, endIP) <= 0; ip = IncrementIP(ip) {
		ipStr := ip.String() + "\n"
		_, err := file.WriteString(ipStr)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return "", err
		}
		time.Sleep(time.Microsecond)
		pgBar.Add(1)
	}
	fmt.Print("\n\n")
	return filePath, nil
}
