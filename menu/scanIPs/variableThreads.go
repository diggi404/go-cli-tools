package menu

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fatih/color"
)

// CheckPorts Main function for scanning the ports received by the user.
func CheckPorts(ip string, ports []string, mutex *sync.Mutex, wg *sync.WaitGroup, timeout *time.Duration, portServices []string, results *[]string, openPorts ...string) {
	defer wg.Done()

	if *timeout == 0 {
		*timeout = time.Second * 10
	}

	start := StartTime()
	for _, port := range ports {
		address := fmt.Sprintf("%s:%s", ip, port)
		conn, err := net.DialTimeout("tcp", address, *timeout)
		if err == nil {
			openPorts = append(openPorts, port)
			serviceInfo, _ := GetServiceInfo(conn)
			sanitizedServiceInfo := SanitizeServiceInfo(serviceInfo)
			portServices = append(portServices, sanitizedServiceInfo)
			conn.Close()
		}
	}
	end := EndTime()
	if len(openPorts) == 0 {
		return
	}
	openPortStr := fmt.Sprintf("%v", openPorts)
	portServiceStr := fmt.Sprintf("%v", portServices)
	green := color.New(color.FgGreen).PrintfFunc()

	// use mutex to lock shared resource for better synchronization between goroutines.
	mutex.Lock()
	defer mutex.Unlock()
	green("%s\t%s\t%s\t%s\t%s\n", ip, openPortStr, portServiceStr, start, end)

	ipData := fmt.Sprintf("%s\t%s\t%s\n", ip, openPortStr, portServiceStr)
	*results = append(*results, ipData)
}
