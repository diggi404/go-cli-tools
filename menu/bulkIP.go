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

func GenIP() (string, error) {
	for attemtps := 0; attemtps < 3; attemtps++ {
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
		totalIps := CountTotalIPs(startIP, endIP)
		bar := pb.Full.Start(totalIps)
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
	fmt.Println("You haved exceeded the try limit!")
	err := errors.New("you haved exceeded the try limit")
	return "", err
}

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

func CountTotalIPs(startIP, endIP net.IP) int {
	start := IpToInt(startIP)
	end := IpToInt(endIP)
	return end - start + 1
}

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
