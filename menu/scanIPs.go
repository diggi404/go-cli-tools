package menu

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
)

// Main function for scanning Bulk IPs.
func ScanIPs(filePath ...string) {
	var (
		user_port      string
		filtered_ports []string
		timeout        int
	)

	// take and filter inputs.
	fmt.Println("Enter the ports you want to scan separated by comma(,). example: 25,587,465")
	fmt.Print(">>>> ")
	fmt.Scanln(&user_port)
	fmt.Print("Enter the timeout in seconds (Default is 1 second) :> ")
	fmt.Scanln(&timeout)

	user_port = strings.TrimSpace(user_port)
	ports := strings.Split(user_port, ",")
	for _, v := range ports {
		if v != "" {
			filtered_ports = append(filtered_ports, v)
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
	var port_services []string
	var results []string
	port_timeout := time.Second * time.Duration(timeout)

	// spawn same number of goroutines as IPs for scanning the ports.
	for _, ip := range ips {
		wg.Add(1)
		go checkPorts(ip, filtered_ports, &mutex, &wg, &port_timeout, port_services, &results)
	}

	// wait for all goroutines to finish...
	wg.Wait()
	writeToFile(results)
	fmt.Println("All checks completed.")
}

// This reads all the IPs from the path passed to it.
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
func checkPorts(ip string, ports []string, mutex *sync.Mutex, wg *sync.WaitGroup, timeout *time.Duration, port_services []string, results *[]string, open_ports ...string) {
	defer wg.Done()

	if *timeout == 0 {
		*timeout = time.Second * 1
	}

	start := time.Now()
	start_hour := start.Hour()
	start_min := start.Minute()
	start_sec := start.Second()

	for _, port := range ports {
		address := fmt.Sprintf("%s:%s", ip, port)
		conn, err := net.DialTimeout("tcp", address, *timeout)
		if err == nil {
			open_ports = append(open_ports, port)
			serviceInfo, _ := getServiceInfo(conn)
			sanitizedServiceInfo := sanitizeServiceInfo(serviceInfo)
			port_services = append(port_services, sanitizedServiceInfo)
			conn.Close()
		}
	}
	if len(open_ports) == 0 {
		return
	}
	open_port_str := fmt.Sprintf("%v", open_ports)
	port_service_str := fmt.Sprintf("%v", port_services)
	green := color.New(color.FgGreen).PrintfFunc()

	end := time.Now()
	end_hour := end.Hour()
	end_min := end.Minute()
	end_sec := end.Second()
	startTime := fmt.Sprintf("%d:%d:%d", start_hour, start_min, start_sec)
	endTime := fmt.Sprintf("%d:%d:%d", end_hour, end_min, end_sec)

	// use mutex to lock shared resource for better synchronization between goroutines.
	mutex.Lock()
	defer mutex.Unlock()
	green("%s\t%s\t%s\t%s\t%s\n", ip, open_port_str, port_service_str, startTime, endTime)

	ip_data := fmt.Sprintf("%s\t%s\t%s\n", ip, open_port_str, port_service_str)
	*results = append(*results, ip_data)
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
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
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
