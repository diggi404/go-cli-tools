package menu

import (
	"errors"
	"fmt"
	"net"
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
	}

	// continuation for both selection from main menu and after generating bulk ips
	ips, err := ReadIPsFromFile(filePath[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	red := color.New(color.FgRed).PrintlnFunc()
	red("IP Address\tOpen Ports\tService     Start\tEnd Time")
	red("----------------------------------------------------------------------------------")

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var portServices []string
	var results []string
	portTimeout := time.Second * time.Duration(timeout)

	if len(ips) > 1000 {

		maxGoroutines := 1000
		chunkSize := len(ips) / maxGoroutines

		if len(ips)%maxGoroutines != 0 {
			chunkSize++
		}
		ipChuncks := make(chan []string, chunkSize)

		for i := 0; i < maxGoroutines; i++ {
			wg.Add(1)
			go checkPorts2(ipChuncks, filteredPorts, &mutex, &wg, &portTimeout, &results)
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
			go checkPorts(ip, filteredPorts, &mutex, &wg, &portTimeout, portServices, &results)
		}
	}

	// wait for all goroutines to finish...
	wg.Wait()
	writeToFile(results)
	fmt.Println("All checks completed.")
}

func checkPorts2(ipChunks <-chan []string, ports []string, mutex *sync.Mutex, wg *sync.WaitGroup, timeout *time.Duration, results *[]string, validIPs ...string) {
	defer wg.Done()

	if *timeout == 0 {
		*timeout = time.Second * 1
	}

	ipsChunk := <-ipChunks

	for _, ip := range ipsChunk {
		var openPorts []string
		var portServices []string

		for _, port := range ports {
			address := fmt.Sprintf("%s:%s", ip, port)
			conn, err := net.DialTimeout("tcp", address, *timeout)
			if err == nil {
				openPorts = append(openPorts, port)
				serviceInfo, _ := getServiceInfo(conn)
				sanitizedServiceInfo := sanitizeServiceInfo(serviceInfo)
				portServices = append(portServices, sanitizedServiceInfo)
				conn.Close()
				// resultFormat := fmt.Sprintf("%s\t%s\t%s\n", ip, port, sanitizedServiceInfo)
				// validIPs = append(validIPs, resultFormat)
			}
		}
		if len(openPorts) != 0 {
			resultFormat := fmt.Sprintf("%s\t%s\t%s\n", ip, openPorts, portServices)
			validIPs = append(validIPs, resultFormat)
		}
	}
	if len(validIPs) == 0 {
		return
	}

	// use mutex to lock shared resource for better synchronization between goroutines.
	mutex.Lock()
	defer mutex.Unlock()

	green := color.New(color.FgGreen).PrintfFunc()
	for _, result := range validIPs {
		green(result)
	}
	*results = append(*results, validIPs...)
}

// func startTime() string {
// 	start := time.Now()
// 	startHour := start.Hour()
// 	startMin := start.Minute()
// 	startSec := start.Second()
// 	startTime := fmt.Sprintf("%d:%d:%d", startHour, startMin, startSec)
// 	return startTime
// }

// func endTime() string {
// 	end := time.Now()
// 	endHour := end.Hour()
// 	endMin := end.Minute()
// 	endSec := end.Second()
// 	endTime := fmt.Sprintf("%d:%d:%d", endHour, endMin, endSec)
// 	return endTime
// }

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

// Main function for scanning the ports received by the user.
func checkPorts(ip string, ports []string, mutex *sync.Mutex, wg *sync.WaitGroup, timeout *time.Duration, portServices []string, results *[]string, openPorts ...string) {
	defer wg.Done()

	start := time.Now()
	startHour := start.Hour()
	startMin := start.Minute()
	startSec := start.Second()

	if *timeout == 0 {
		*timeout = time.Second * 1
	}

	for _, port := range ports {
		address := fmt.Sprintf("%s:%s", ip, port)
		conn, err := net.DialTimeout("tcp", address, *timeout)
		if err == nil {
			openPorts = append(openPorts, port)
			serviceInfo, _ := getServiceInfo(conn)
			sanitizedServiceInfo := sanitizeServiceInfo(serviceInfo)
			portServices = append(portServices, sanitizedServiceInfo)
			conn.Close()
		}
	}
	if len(openPorts) == 0 {
		return
	}
	openPortStr := fmt.Sprintf("%v", openPorts)
	portServiceStr := fmt.Sprintf("%v", portServices)
	green := color.New(color.FgGreen).PrintfFunc()

	end := time.Now()
	endHour := end.Hour()
	endMin := end.Minute()
	endSec := end.Second()
	startTime := fmt.Sprintf("%d:%d:%d", startHour, startMin, startSec)
	endTime := fmt.Sprintf("%d:%d:%d", endHour, endMin, endSec)

	// use mutex to lock shared resource for better synchronization between goroutines.
	mutex.Lock()
	defer mutex.Unlock()
	green("%s\t%s\t%s\t%s\t%s\n", ip, openPortStr, portServiceStr, startTime, endTime)

	ipData := fmt.Sprintf("%s\t%s\t%s\n", ip, openPortStr, portServiceStr)
	*results = append(*results, ipData)
}

func writeToFile(results []string) {
	dirPath, err := setupDir()
	if err != nil {
		return
	}
	filePath := dirPath + "/scanned_ips.txt"
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	for _, result := range results {
		_, err = file.WriteString(result)
		if err != nil {
			return
		}
	}
}

// clear the port service info off any newline characters.
func sanitizeServiceInfo(serviceInfo string) string {
	sanitized := strings.ReplaceAll(serviceInfo, "\n", " ")
	sanitized = strings.TrimSpace(sanitized)
	return sanitized
}

// get the services running on each port connection created.
func getServiceInfo(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return "Timeout", err
		} else {
			return "Timeout", err
		}
	}
	data := string(buffer[:n])
	return data, nil
}

// initialize the default directory path for keeping results of scanned IPs
func setupDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	dirPath := cwd + "/ip_scans"
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
