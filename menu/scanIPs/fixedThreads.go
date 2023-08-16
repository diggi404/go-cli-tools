package menu

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fatih/color"
)

func CheckPorts2(ipChunks <-chan []string, ports []string, mutex *sync.Mutex, wg *sync.WaitGroup, timeout *time.Duration, results *[]string, validIPs ...string) {
	defer wg.Done()

	if *timeout == 0 {
		*timeout = time.Second * 10
	}

	// get pushed data to channel at a go.
	ipsChunk := <-ipChunks

	for _, ip := range ipsChunk {
		var openPorts []string
		var portServices []string
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
		if len(openPorts) != 0 {
			resultFormat := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", ip, openPorts, portServices, start, end)
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
