package menu

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func GenIP() {
	for attemtps := 0; attemtps < 3; attemtps++ {
		file, err := os.Create("generated_ips.txt")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
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

		for ip := startIP; ip != nil && bytesCompare(ip, endIP) <= 0; ip = incrementIP(ip) {
			ipStr := ip.String() + "\n"
			_, err := file.WriteString(ipStr)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			time.Sleep(time.Millisecond)
			bar.Increment()
		}
		bar.Finish()
		return
	}
	fmt.Println("You haved exceeded the try limit!")
	os.Exit(1)
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
