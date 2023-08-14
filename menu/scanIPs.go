package menu

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/ncruces/zenity"
	"github.com/olekukonko/tablewriter"
)

func ScanIPs(table *tablewriter.Table, filePath ...string) {
	var (
		user_port      string
		filtered_ports []string
	)
	fmt.Println("Enter the ports you want to scan separated by comma(,). example: 25,587,465")
	fmt.Print(">>>> ")
	fmt.Scanln(&user_port)
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
		fmt.Printf("fileName: %v\n", fileName)
		filePath = append(filePath, fileName)
	}
	ips, err := readIPsFromFile(filePath[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var wg sync.WaitGroup
	var tableMutex sync.Mutex

	for _, ip := range ips {
		wg.Add(1)
		go checkPorts(ip, filtered_ports, table, &tableMutex, &wg)
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

func checkPorts(ip string, ports []string, table *tablewriter.Table, tableMutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	var open_ports []string
	for _, port := range ports {
		address := fmt.Sprintf("%s:%s", ip, port)
		conn, err := net.Dial("tcp", address)
		if err == nil {
			open_ports = append(open_ports, port)
			conn.Close()
		}
	}

	tableMutex.Lock()
	defer tableMutex.Unlock()
	openPortsStr := fmt.Sprintf("%v", open_ports)
	table.Append([]string{ip, openPortsStr})
	table.Render()
}
