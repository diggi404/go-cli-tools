package cpanel

import (
	"bufio"
	"fmt"
	"go_cli/fileutil"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
)

func CPanelChecker() {
	blue := color.New(color.FgHiBlue).PrintFunc()
	red := color.New(color.FgRed).PrintfFunc()
	reader := bufio.NewReader(os.Stdin)
	red("\nYour wordlist should be in this format > https://example.com:2083|username|password\n")
	red("NOTE: Credentials with invalid format will be skipped or ignored automatically!\n")
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
		return
	}
	wordlist, err := fileutil.ReadFromFile(filePath)
	if err != nil {
		red("err: %v\n", err)
		return
	}
	color.New(color.FgHiMagenta).Printf("\nTotal Credentials: %v\n", len(wordlist))

	maxWorkers := 1000
	chunkSize := len(wordlist) / maxWorkers

	if len(wordlist)%maxWorkers != 0 {
		chunkSize++
	}

	wordlistChunks := make(chan []string, chunkSize)

	var mutex sync.Mutex
	var wg sync.WaitGroup
	totalChecks := 0

	currentTime := time.Now().Unix()
	fileName := fmt.Sprintf("hits_%v.txt", currentTime)
	file, err := fileutil.WriteToFile("cpanel_logs", fileName)
	if err != nil {
		red("err: %v\n", err)
		return
	}
	fmt.Println()

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go ProcessCreds(wordlistChunks, &wg, &mutex, file, &totalChecks)
	}

	for i := 0; i < len(wordlist); i += chunkSize {
		end := i + chunkSize
		if end > len(wordlist) {
			end = len(wordlist)
		}
		wordlistChunks <- wordlist[i:end]
	}
	close(wordlistChunks)
	wg.Wait()
	color.New(color.FgMagenta).Println("all done. Thanks for using the tool.")
}
