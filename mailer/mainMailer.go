package mailer

import (
	"bufio"
	"crypto/tls"
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

type MailOut struct {
	Subject    string
	FromEmail  string
	FromName   string
	Message    string
	IsMsgPlain bool
}

type MailOutResults struct {
	Success int
	Fails   int
}

func Mailer() {
	var (
		filteredCreds []string
	)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your SMTP Credentials. Format >  HOST,PORT,USERNAME,PASSWORD")
	fmt.Print(">>> ")
	smtpCreds, _ := reader.ReadString('\n')
	splitedCreds := strings.Split(smtpCreds, ",")

	for _, creds := range splitedCreds {
		trimedCreds := strings.TrimSpace(creds)
		if len(trimedCreds) != 0 {
			filteredCreds = append(filteredCreds, trimedCreds)
		}
	}

	fmt.Println("verifying SMTP credentials...")

	host, portStr, username, password := filteredCreds[0], filteredCreds[1], filteredCreds[2], filteredCreds[3]
	port, _ := strconv.Atoi(portStr)
	dailer := gomail.NewDialer(host, port, username, password)
	dailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	smtpConn, err := dailer.Dial()
	if err != nil {
		errMsg := fmt.Sprintf("err: %v\n", err)
		color.HiRed(errMsg)
		return
	}
	defer smtpConn.Close()

	color.New(color.FgGreen).Println("SMTP connection was successful.")
	mailOpts := MailOut{FromEmail: username}

	fmt.Println("Select your email list: ")
	filePath, err := zenity.SelectFile(
		zenity.FileFilters{
			{Patterns: []string{"*.txt"}, CaseFold: false},
		})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	emailList, err := fileutil.ReadFromFile(filePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("Total Emails: %v\n", len(emailList))

	var rawMsgType string
	fmt.Print("What type of content are you sending? plain/html :> ")
	fmt.Scanln(&rawMsgType)
	msgType := strings.ToLower(strings.TrimSpace(rawMsgType))

	if msgType == "html" || strings.Contains(msgType, "html") {

		fmt.Println("Select your html letter: ")
		filePath, err := zenity.SelectFile(
			zenity.FileFilters{
				{Patterns: []string{"*.html"}, CaseFold: false},
			})
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		htmlByte, err := os.ReadFile(filePath)

		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		mailOpts.Message = string(htmlByte)
		mailOpts.IsMsgPlain = false
	} else {

		fmt.Print("Enter your Message :> ")
		msg, _ := reader.ReadString('\n')
		mailOpts.Message = msg
		mailOpts.IsMsgPlain = true
	}

	fmt.Print("Enter Message Subject :> ")
	msgSubject, _ := reader.ReadString('\n')
	mailOpts.Subject = strings.TrimSpace(msgSubject)

	fmt.Print("Enter from name :> ")
	fromName, _ := reader.ReadString('\n')
	mailOpts.FromName = strings.TrimSpace(fromName)

	maxWorkers := 1000
	chunkSize := len(emailList) / maxWorkers

	if len(emailList)%maxWorkers != 0 {
		chunkSize++
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	bar := progressbar.NewOptions(len(emailList),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("Sending emails..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	results := MailOutResults{}

	emailListChunks := make(chan []string, chunkSize)

	msgOpts := gomail.NewMessage()
	msgOpts.SetAddressHeader("From", mailOpts.FromEmail, mailOpts.FromName)
	msgOpts.SetHeader("Subject", mailOpts.Subject)
	if mailOpts.IsMsgPlain {
		msgOpts.SetBody("text/plain", mailOpts.Message)
	} else {
		msgOpts.SetBody("text/html", mailOpts.Message)
	}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go SendMail(emailListChunks, &wg, &mutex, mailOpts, &smtpConn, msgOpts, bar, &results)
	}

	for i := 0; i < len(emailList); i += chunkSize {
		end := i + chunkSize
		if end > len(emailList) {
			end = len(emailList)
		}
		emailListChunks <- emailList[i:end]
	}
	close(emailListChunks)
	wg.Wait()
	successMsg := fmt.Sprintf("\nSent: %d", results.Success)
	failMsg := fmt.Sprintf("Fails: %d", results.Fails)
	color.New(color.FgGreen).Println(successMsg)
	color.New(color.FgRed).Println(failMsg)
}
