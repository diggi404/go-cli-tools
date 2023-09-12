package scanips

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

func CheckPorts2(ipChunks <-chan []string, ports []string, mutex *sync.Mutex, wg *sync.WaitGroup, timeout time.Duration, files []*os.File, totalChecks *int) {
	defer wg.Done()

	if timeout == 0 {
		timeout = time.Second * 10
	}

	// get pushed data to channel at a go.
	ipsChunk := <-ipChunks

	for _, ip := range ipsChunk {
		var openPorts []string
		var portServices []string
		for i, port := range ports {
			address := fmt.Sprintf("%s:%s", ip, port)
			conn, err := net.DialTimeout("tcp", address, timeout)
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
			}
		}
		mutex.Lock()
		if len(openPorts) != 0 {
			*totalChecks++
			openPortStr := fmt.Sprintf("%v", openPorts)
			portServiceStr := fmt.Sprintf("%v", portServices)
			color.New(color.FgBlue).Printf("%d: -> ", *totalChecks)
			color.New(color.FgGreen).Printf("%s\t%s\t%s\n", ip, openPortStr, portServiceStr)
		} else {
			*totalChecks++
		}
		mutex.Unlock()
	}
}
