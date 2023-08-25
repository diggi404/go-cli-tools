package scanips

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

func CheckPorts2(ipChunks <-chan []string, ports []string, mutex *sync.Mutex, wg *sync.WaitGroup, timeout *time.Duration, file *os.File, totalChecks *int) {
	defer wg.Done()

	if *timeout == 0 {
		*timeout = time.Second * 10
	}

	green := color.New(color.FgGreen).PrintfFunc()
	blue := color.New(color.FgBlue).PrintfFunc()

	// get pushed data to channel at a go.
	ipsChunk := <-ipChunks

	for _, ip := range ipsChunk {
		var openPorts []string
		var portServices []string
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
		if len(openPorts) != 0 {
			mutex.Lock()
			*totalChecks += 1
			blue("%d: -> ", *totalChecks)
			results := fmt.Sprintf("%s\t%s\t%s\n", ip, openPorts, portServices)
			green(results)
			file.WriteString(results)
			mutex.Unlock()
		} else {
			mutex.Lock()
			*totalChecks += 1
			mutex.Unlock()
		}
	}
}
