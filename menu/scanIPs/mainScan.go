package menu

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
)

// ScanIPs Main function for scanning Bulk IPs.
func ScanIPs(filePath ...string) {
	var (
		userPort      string
		filteredPorts []string
		timeout       int
	)

	// take and filter inputs.
	fmt.Println("Enter the ports you want to scan separated by comma(,). example: 25,587,465")
	fmt.Print(">>>> ")
	fmt.Scanln(&userPort)
	fmt.Print("Enter the timeout in seconds (Default is 1 second) :> ")
	fmt.Scanln(&timeout)

	userPort = strings.TrimSpace(userPort)
	ports := strings.Split(userPort, ",")
	for _, v := range ports {
		if v != "" {
			filteredPorts = append(filteredPorts, v)
		}
	}

	// handles direct selection from Main Menu
	if len(filePath) == 0 {
		fmt.Println("Select your file: ")
		fileName, err := zenity.SelectFile(
			zenity.FileFilters{
				{Name: "IP list", Patterns: []string{"*.txt"}, CaseFold: false},
			})
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		filePath = append(filePath, fileName)
		fmt.Printf("filePath: %v\n", fileName)
	}

	// continuation for both selection from main menu and after generating bulk ips
	ips, err := ReadIPsFromFile(filePath[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Total IPs: %v\n", len(ips))
	red := color.New(color.FgRed).PrintlnFunc()
	red("IP Address\tOpen Ports\tService     Start\tEnd Time")
	red("----------------------------------------------------------------------------------")

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var portServices []string
	var results []string
	portTimeout := time.Second * time.Duration(timeout)

	if len(ips) > 1000 {

		maxWorkers := 1000
		chunkSize := len(ips) / maxWorkers

		if len(ips)%maxWorkers != 0 {
			chunkSize++
		}
		ipChuncks := make(chan []string, chunkSize)

		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go CheckPorts2(ipChuncks, filteredPorts, &mutex, &wg, &portTimeout, &results)
		}

		for i := 0; i < len(ips); i += chunkSize {
			end := i + chunkSize
			if end > len(ips) {
				end = len(ips)
			}
			ipChuncks <- ips[i:end]
		}

		close(ipChuncks)
	} else {
		// spawn same number of goroutines as IPs for scanning the ports.
		for _, ip := range ips {
			wg.Add(1)
			go CheckPorts(ip, filteredPorts, &mutex, &wg, &portTimeout, portServices, &results)
		}
	}

	// wait for all goroutines to finish...
	wg.Wait()
	WriteToFile(results)
	fmt.Println("All checks completed.")
}
