package mailer

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/gomail.v2"
)

func HandleGmailSMTP(emailChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, smtpConn *SmtpConnOpts, smtpCreds SmtpOpts, msgOpts *gomail.Message, pgBar *progressbar.ProgressBar) {
	defer wg.Done()
	emailList := <-emailChunks
	for _, email := range emailList {
		mutex.Lock()
		msgOpts.SetHeader("To", email)
		for {
			err := gomail.Send(smtpConn.Conn, msgOpts)
			if err == nil {
				pgBar.Add(1)
				smtpConn.NewConn = false
				break
			} else if strings.Contains(err.Error(), "SMTP Daily user sending quota exceeded.") {
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
