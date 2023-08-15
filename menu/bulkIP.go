package menu

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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
		dirPath, err := checkDir()
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
		for ip := startIP; ip != nil && bytesCompare(ip, endIP) <= 0; ip = incrementIP(ip) {
			ipStr := ip.String() + "\n"
			_, err := file.WriteString(ipStr)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return "", err
			}
			time.Sleep(time.Millisecond)
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

// This function checks if the default directory of the application is available and if not it makes a new dir.
func checkDir() (string, error) {
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

// CountTotalIPs Calculate the total number of IPs that can be generated within the given range;
// This is used or the cli progress bar.
func CountTotalIPs(startIP, endIP net.IP) int {
	start := IpToInt(startIP)
	end := IpToInt(endIP)
	return end - start + 1
}

// IpToInt This will convert an IP address[byte] into it's full integer value;
// Used in the CountTotalIPs function to calc the number of IPs within a range.
func IpToInt(ip net.IP) int {
	octets := strings.Split(ip.String(), ".")
	if len(octets) != 4 {
		return 0
	}

	var result int
	for _, octetStr := range octets {
		octet, err := strconv.Atoi(octetStr)
		if err != nil {
			return 0
		}
		result = result*256 + octet
	}

	return result
}

// Compares the starting and ending IPs. Used as a condition in the for loop generating the IPs.
func bytesCompare(a, b net.IP) int {
	for i := range a {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

// This increments the IPs and which then becomes the starting IP for the next iteration.
func incrementIP(ip net.IP) net.IP {
	nextIP := make(net.IP, len(ip))
	copy(nextIP, ip)
	for i := len(nextIP) - 1; i >= 0; i-- {
		nextIP[i]++
		if nextIP[i] > 0 {
			break
		}
	}
	return nextIP
}
