package scanips

import (
	"bufio"
	"fmt"
	"go_cli/fileutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
)

func CloseFile(files []*os.File) {
	for _, file := range files {
		file.Close()
	}
}

// ScanIPs Main function for scanning Bulk IPs.
func ScanIPs(filePath ...string) {
	var (
		filteredPorts []string
		timeout       int
	)
	blue := color.New(color.FgHiBlue).PrintFunc()
	red := color.New(color.FgRed).PrintfFunc()
	// take and filter inputs.
	reader := bufio.NewReader(os.Stdin)
	blue("\nEnter the ports you want to scan separated by comma(,). example: 22,3389,2083\n")
	blue(">>>> ")
	userPort, err := reader.ReadString('\n')
	if err != nil {
		red("err: %v\n", err)
		return
	} else if userPort == "\n" {
		red("invalid input. Exiting Program...\n")
		return
	}
	blue("\nEnter the timeout in seconds (Default = 10s) :> ")
	fmt.Scanln(&timeout)

	// filter all entered ports
	ports := strings.Split(userPort, ",")
	for _, port := range ports {
		portTrimmed := strings.TrimSpace(port)
		if len(portTrimmed) != 0 {

			filteredPorts = append(filteredPorts, portTrimmed)
		}
	}

	// handles direct selection from Main Menu
	if len(filePath) == 0 {
		blue("\nSelect your file: \n")
		fileName, err := zenity.SelectFile(
			zenity.FileFilters{
				{Patterns: []string{"*.txt"}, CaseFold: false},
			})
		if err != nil {
			red("err: %v\n", err)
			return
		}
		filePath = append(filePath, fileName)
	}

	// continuation for both selection from main menu and after generating bulk ips
	ips, err := fileutil.ReadFromFile(filePath[0])
	if err != nil {
		red("err: %v\n", err)
		return
	}
	color.New(color.FgHiMagenta).Printf("\nTotal IPs: %v\n", len(ips))
	red("\nIP Address\tOpen Ports\tService\n")
	red("----------------------------------------------------------------------------------\n")

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var portServices []string
	totalChecks := 0
	portTimeout := time.Second * time.Duration(timeout)

	curentTime := time.Now().Unix()
	dirName := fmt.Sprintf("ip_scans/scanned_ips_%v", curentTime)
	var files []*os.File
	for _, port := range filteredPorts {
		file, err := fileutil.WriteToFile(dirName, port+".txt")
		if err != nil {
			red("err: %v\n", err)
			return
		}
		files = append(files, file)
	}

	defer CloseFile(files)

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
			go CheckPorts2(ipChunks, filteredPorts, &mutex, &wg, portTimeout, files, &totalChecks)
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
			go CheckPorts(ip, filteredPorts, &mutex, &wg, &portTimeout, portServices, files, &totalChecks)
		}
	}

	// wait for all goroutines to finish...
	wg.Wait()
	color.New(color.FgMagenta).Println("\nAll checks completed. Thanks for using the tool.")
}
