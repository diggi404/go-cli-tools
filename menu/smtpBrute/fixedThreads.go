package smtpbrute

import (
	"fmt"
	"sync"
)

func ProcessCredentials(wordList <-chan []string, testEmail string, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	wordListChunks := <-wordList
	fmt.Printf("wordListChunks: %v\n", wordListChunks)
	for _, creds := range wordListChunks {
		splitedCreds, err := FilterGmailCreds(creds)
		if err != nil {
			ConnectSMTP(splitedCreds, testEmail)
		} else {
			smtpCreds, err := LookupDomain(splitedCreds)
			if err == nil {
				ConnectSMTP(smtpCreds, testEmail)
			}
		}
	}
}
