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
func CheckPorts(ip string, ports []string, mutex *sync.Mutex, wg *sync.WaitGroup, timeout *time.Duration, portServices []string, files []*os.File, totalChecks *int, openPorts ...string) {
	defer wg.Done()

	if *timeout == 0 {
		*timeout = time.Second * 10
	}

	for i, port := range ports {
		address := fmt.Sprintf("%s:%s", ip, port)
		conn, err := net.DialTimeout("tcp", address, *timeout)
		if err == nil {
			openPorts = append(openPorts, port)
			serviceInfo, _ := GetServiceInfo(conn)
			sanitizedServiceInfo := SanitizeServiceInfo(serviceInfo)
			portServices = append(portServices, sanitizedServiceInfo)
			conn.Close()
			file := files[i]
			result := fmt.Sprintf("%s\t\t[%s]\n", ip, sanitizedServiceInfo)
			mutex.Lock()
			file.WriteString(result)
			mutex.Unlock()
		} else {
			mutex.Lock()
			*totalChecks += 1
			mutex.Unlock()
		}
	}
	if len(openPorts) == 0 {
		return
	}

	// use mutex to lock shared resource for better synchronization between goroutines.
	mutex.Lock()
	*totalChecks += 1
	openPortStr := fmt.Sprintf("%v", openPorts)
	portServiceStr := fmt.Sprintf("%v", portServices)
	color.New(color.FgBlue).Printf("%d: -> ", *totalChecks)
	color.New(color.FgGreen).Printf("%s\t%s\t%s\n", ip, openPortStr, portServiceStr)
	mutex.Unlock()
}
