package menu

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ncruces/zenity"
)

func ScanIPs(filePath ...string) {
	var (
		user_port      string
		filtered_ports []string
		timeout        time.Duration
	)
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
	ips, err := readIPsFromFile(filePath[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("IP Address\tOpen Ports")
	fmt.Println("-------------------------------------")

	var wg sync.WaitGroup

	for _, ip := range ips {
		wg.Add(1)
		go checkPorts(ip, filtered_ports, &wg, timeout)
	}

	wg.Wait()
	fmt.Println("All checks completed.")
}

func readIPsFromFile(fileName string) ([]string, error) {
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

func checkPorts(ip string, ports []string, wg *sync.WaitGroup, timeout time.Duration, open_ports ...string) {
	defer wg.Done()
	if timeout == 0 {
		timeout = 1
	}
	for _, port := range ports {
		address := fmt.Sprintf("%s:%s", ip, port)
		conn, err := net.DialTimeout("tcp", address, time.Second*timeout)
		if err == nil {
			open_ports = append(open_ports, port)
			conn.Close()
		}
	}
	if len(open_ports) == 0 {
		return
	}
	openPortsStr := fmt.Sprintf("%v", open_ports)
	fmt.Printf("%s\t%s\n", ip, openPortsStr)
	dirPath, err := setupDir()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	filePath := dirPath + "/scanned_ips.txt"
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	ip_data := fmt.Sprintf("%s\t%s\n", ip, openPortsStr)
	_, err = file.WriteString(ip_data)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

}

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
