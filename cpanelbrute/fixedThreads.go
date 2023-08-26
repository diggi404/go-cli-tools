package cpanelbrute

import (
	"sync"

	"github.com/olekukonko/tablewriter"
)

func HandleBrute(target string, wordlistChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, table *tablewriter.Table) {
	defer wg.Done()

	credChunks := <-wordlistChunks

	for _, creds := range credChunks {
		validCreds, err := MakeRequests(target, creds)
		if err == nil {
			username, password, url := validCreds[0], validCreds[1], validCreds[2]
			mutex.Lock()
			table.Append([]string{url, username, password})
			table.Render()
			mutex.Unlock()
		}
	}

}
