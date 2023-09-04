package scanips

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// CheckPorts Main function for scanning the ports received by the user.
func CheckPorts(ip string, ports []string, mutex *sync.Mutex, wg *sync.WaitGroup, timeout *time.Duration, portServices []string, file *os.File, totalChecks *int, openPorts ...string) {
	defer wg.Done()

	if *timeout == 0 {
		*timeout = time.Second * 10
	}

	for _, port := range ports {
		address := fmt.Sprintf("%s:%s", ip, port)
		conn, err := net.DialTimeout("tcp", address, *timeout)
		if err == nil {
			openPorts = append(openPorts, port)
			serviceInfo, _ := GetServiceInfo(conn)
			sanitizedServiceInfo := SanitizeServiceInfo(serviceInfo)
			portServices = append(portServices, sanitizedServiceInfo)
			conn.Close()
		} else {
			mutex.Lock()
			*totalChecks += 1
			mutex.Unlock()
		}
	}
	if len(openPorts) == 0 {
		return
	}
	openPortStr := fmt.Sprintf("%v", openPorts)
	portServiceStr := fmt.Sprintf("%v", portServices)
	green := color.New(color.FgGreen).PrintfFunc()
	blue := color.New(color.FgBlue).PrintfFunc()

	// use mutex to lock shared resource for better synchronization between goroutines.
	mutex.Lock()
	defer mutex.Unlock()
	*totalChecks += 1
	blue("%d: -> ", *totalChecks)
	results := fmt.Sprintf("%s\t%s\t%s\n", ip, openPortStr, portServiceStr)
	green(results)
	file.WriteString(results)
}
