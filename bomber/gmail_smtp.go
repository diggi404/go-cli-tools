package bomber

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/gomail.v2"
)

func HandleGmailSMTP(articleChunks <-chan Article, workingChan <-chan []int, wg *sync.WaitGroup, mutex *sync.Mutex, msgOptions *gomail.Message, smtpCreds SmtpOpts, pgBar *progressbar.ProgressBar, smtpConn *SmtpConnOpts) {
	defer wg.Done()
	article := <-articleChunks
	rounds := <-workingChan
	for range rounds {
		if article.Author == "" {
			article.Author = "Fabrizio Romano"
		}
		rand.New(rand.NewSource(time.Now().UnixNano()))
		randomInt := rand.Intn(1000)
		randNum := strconv.Itoa(randomInt)
		mutex.Lock()
		msgOptions.SetAddressHeader("From", smtpCreds.Username, article.Author)
		msgOptions.SetHeader("Subject", article.Title+" "+randNum)
		msgOptions.SetBody("text/plain", article.Description)
		for {
			err := gomail.Send(smtpConn.Conn, msgOptions)
			if err == nil {
				pgBar.Add(1)
				smtpConn.NewConn = false
				break
			} else if strings.Contains(err.Error(), "Daily user sending quota exceeded.") {
				color.New(color.FgRed).Printf("\nDaily limit reached. Try again after 24hrs or use a different account.\n")
				os.Exit(1)
			} else if smtpConn.NewConn && smtpConn.NumErrors >= 2 {
				color.New(color.FgRed).Printf("\nThe SMTP has been rate limited. Try again later!\n")
				os.Exit(1)
			} else if err != nil {
				if smtpConn.NewConn {
					smtpConn.NumErrors++
				}
				color.New(color.FgRed).Printf("\nerr: %v\n", err)
				color.New(color.FgGreen).Printf("\nretrying in 3 mins...\n")

				duration := 3 * time.Minute
				startTime := time.Now()
				for {
					currentTime := time.Now()
					remainingTime := duration - currentTime.Sub(startTime)
					if remainingTime <= 0 {
						break
					}
					color.New(color.FgHiMagenta).Printf("\rTime Remaining: %02d:%02d",
						int(remainingTime.Minutes()),
						int(remainingTime.Seconds())%60)
					time.Sleep(time.Second)
				}
				color.New(color.FgGreen).Printf("\n\nusing a different Connection\n")
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
