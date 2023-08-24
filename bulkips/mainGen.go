package bulkips

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

// GenIP main function to generate the bulk IPs.
func GenIP() (string, error) {
	for attempts := 0; attempts < 3; attempts++ {
		var ip1, ip2 string
		fmt.Print("Enter the starting IP: ")
		fmt.Scanln(&ip1)
		startIP := net.ParseIP(ip1).To4()
		if startIP == nil {
			fmt.Println("Please enter a valid IP address.")
			continue
		}
		fmt.Print("Enter the ending IP: ")
		fmt.Scanln(&ip2)
		endIP := net.ParseIP(ip2).To4()
		if endIP == nil {
			fmt.Println("Please enter a valid IP address.")
			continue
		}

		// initialize the bar cli output
		totalIps := CountTotalIPs(startIP, endIP)
		bar := pb.Full.Start(totalIps)

		// prepare the result file's directory and file.
		dirPath, err := CheckDir()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return "", err
		}
		filePath := dirPath + "/generated_ips.txt"
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
			bar.Increment()
		}
		bar.Finish()
		fmt.Print("\n")
		return filePath, nil
	}
	fmt.Println("You have exceeded the try limit!")
	err := errors.New("you have exceeded the try limit")
	return "", err
}
