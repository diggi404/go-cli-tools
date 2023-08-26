package scanips

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
)

// ScanIPs Main function for scanning Bulk IPs.
func ScanIPs(filePath ...string) {
	var (
		filteredPorts []string
		timeout       int
	)

	// take and filter inputs.
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the ports you want to scan separated by comma(,). example: 22,3389,2083")
	fmt.Print(">>>> ")
	userPort, _ := reader.ReadString('\n')
	fmt.Print("Enter the timeout in seconds (Default = 10s) :> ")
	fmt.Scanln(&timeout)

	// filter all entered ports
	ports := strings.Split(userPort, ",")
	for _, port := range ports {
		portTrimed := strings.TrimSpace(port)
		if len(portTrimed) != 0 {

			filteredPorts = append(filteredPorts, portTrimed)
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
	red("IP Address\tOpen Ports\tService")
	red("----------------------------------------------------------------------------------")

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var portServices []string
	var totalChecks int
	portTimeout := time.Second * time.Duration(timeout)

	file := WriteToFile()
	defer file.Close()

	// spawn a fixed number of goroutines for files contain more than 1000 IPs
	if len(ips) > 1000 {

		maxWorkers := 1000
		chunkSize := len(ips) / maxWorkers

		if len(ips)%maxWorkers != 0 {
			chunkSize++
		}
		ipChunks := make(chan []string, chunkSize)

		// spawn goroutines which will be reading data from the ipChunks channel concurrently.
		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go CheckPorts2(ipChunks, filteredPorts, &mutex, &wg, &portTimeout, &file, &totalChecks)
		}

		// share IPs among goroutines by sending calculated chunk data size to worker channel.
		for i := 0; i < len(ips); i += chunkSize {
			end := i + chunkSize
			if end > len(ips) {
				end = len(ips)
			}
			ipChunks <- ips[i:end]
		}

		close(ipChunks)
	} else {
		// spawn same number of goroutines as IPs for scanning the ports.
		for _, ip := range ips {
			wg.Add(1)
			go CheckPorts(ip, filteredPorts, &mutex, &wg, &portTimeout, portServices, &file, &totalChecks)
		}
	}

	// wait for all goroutines to finish...
	wg.Wait()
	fmt.Println("All checks completed.")
}
