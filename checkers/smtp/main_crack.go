package smtp

import (
	"bufio"
	"fmt"
	"go_cli/fileutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
)

func SMTPChecker() {
	blue := color.New(color.FgHiBlue).PrintFunc()
	red := color.New(color.FgRed).PrintfFunc()
	reader := bufio.NewReader(os.Stdin)
	var testEmail string
	blue("\nEnter test email :> ")
	fmt.Scanln(&testEmail)
	if len(testEmail) == 0 {
		red("invalid input. Exiting Program...\n")
		return
	}
	blue("\nPress Enter to select your wordlist: ")
	_, err := reader.ReadString('\n')
	if err != nil {
		red("err: %v\n", err)
		return
	}
	filePath, err := zenity.SelectFile(
		zenity.FileFilters{
			{Patterns: []string{"*.txt"}, CaseFold: false},
		})
	if err != nil {
		red("err: %v\n", err)
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
		red("err: %v\n", err)
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
	color.New(color.FgMagenta).Println("\nall done. Thanks for using the tool.")
}
