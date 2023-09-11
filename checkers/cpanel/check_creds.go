package cpanel

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func ProcessCreds(target string, wordlistChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, table *tablewriter.Table, file *os.File, totalChecks *int) {
	defer wg.Done()

	credChunks := <-wordlistChunks
	for _, creds := range credChunks {
		validCreds, err := MakeRequest(target, creds)
		mutex.Lock()
		if err == nil {
			username, password, url := validCreds[0], validCreds[1], validCreds[2]
			table.Append([]string{url, username, password})
			table.Render()
			result := fmt.Sprintf("%s => %s", target, creds)
			file.WriteString(result)
			os.Exit(0)
		} else {
			*totalChecks++
			color.New(color.FgBlue).Printf("%d: -> ", *totalChecks)
			color.New(color.FgRed).Printf("%s\n", creds)
		}
		mutex.Unlock()
	}

}
