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

type SmtpOpts struct {
	Host     string
	Port     string
	Username string
	Password string
}

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

type SmtpConnOpts struct {
	Conn      gomail.SendCloser
	NumErrors int
	NewConn   bool
}

func SmtpConnClose(conns []gomail.SendCloser) {
	for _, conn := range conns {
		if conn != nil {
			conn.Close()
		}
	}
}

func Mailer() {
	var (
		filteredCreds []string
	)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your SMTP Credentials. Format >  HOST,PORT,USERNAME,PASSWORD")
	fmt.Print(">>> ")
	smtpCredsStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	} else if smtpCredsStr == "\n" {
		fmt.Println("invalid input!")
		return
	}
	splittedCreds := strings.Split(smtpCredsStr, ",")

	for _, creds := range splittedCreds {
		trimedCreds := strings.TrimSpace(creds)
		if len(trimedCreds) != 0 {
			filteredCreds = append(filteredCreds, trimedCreds)
		}
	}

	fmt.Println("\nverifying SMTP credentials...")

	host, portStr, username, password := filteredCreds[0], filteredCreds[1], filteredCreds[2], filteredCreds[3]
	port, _ := strconv.Atoi(portStr)
	dailer := gomail.NewDialer(host, port, username, password)
	dailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	conn, err := dailer.Dial()
	if err != nil {
		errMsg := fmt.Sprintf("err: %v\n", err)
		color.HiRed(errMsg)
		return
	}
	smtpCreds := SmtpOpts{
		Host:     host,
		Port:     portStr,
		Username: username,
		Password: password,
	}
	smtpConn := SmtpConnOpts{Conn: conn, NewConn: true}
	defer smtpConn.Conn.Close()

	color.New(color.FgGreen).Printf("\nSMTP connection established successfully :)\n")
	mailOpts := MailOut{FromEmail: username}

	fmt.Println("\nSelect your email list: ")
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
	color.New(color.FgHiMagenta).Printf("\nTotal Emails: %v\n", len(emailList))

	var rawMsgType string
	fmt.Print("\nWhat type of content are you sending? plain/html :> ")
	fmt.Scanln(&rawMsgType)
	if len(rawMsgType) == 0 {
		fmt.Println("invalid choice!")
		return
	}
	msgType := strings.ToLower(strings.TrimSpace(rawMsgType))

	if msgType == "html" || strings.Contains(msgType, "html") {

		fmt.Println("\nSelect your html letter: ")
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

		fmt.Print("\nEnter your Message :> ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		} else if msg == "\n" {
			fmt.Println("invalid input!")
			return
		}
		mailOpts.Message = msg
		mailOpts.IsMsgPlain = true
	}

	fmt.Print("\nEnter Message Subject :> ")
	msgSubject, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	} else if msgSubject == "\n" {
		fmt.Println("invalid input!")
		return
	}
	mailOpts.Subject = strings.TrimSpace(msgSubject)

	fmt.Print("\nEnter from name :> ")
	fromName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	} else if fromName == "\n" {
		fmt.Println("invalid input!")
		return
	}
	fmt.Println()
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

	emailListChunks := make(chan []string, chunkSize)

	msgOpts := gomail.NewMessage()
	msgOpts.SetAddressHeader("From", mailOpts.FromEmail, mailOpts.FromName)
	msgOpts.SetHeader("Subject", mailOpts.Subject)
	if mailOpts.IsMsgPlain {
		msgOpts.SetBody("text/plain", mailOpts.Message)
	} else {
		msgOpts.SetBody("text/html", mailOpts.Message)
	}

	if smtpCreds.Host == "smtp.gmail.com" {
		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go HandleGmailSMTP(emailListChunks, &wg, &mutex, &smtpConn, smtpCreds, msgOpts, bar)
		}
	} else {
		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go SendMail(emailListChunks, &wg, &mutex, &smtpConn, smtpCreds, msgOpts, bar)
		}
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
	fmt.Println("\n\nall done.")
}
