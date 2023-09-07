package mailer

import (
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/gomail.v2"
)

func SendMail(emailChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, smtpConn *SmtpConnOpts, smtpCreds SmtpOpts, msgOpts *gomail.Message, pgBar *progressbar.ProgressBar) {
	defer wg.Done()
	emailList := <-emailChunks
	for _, email := range emailList {
		msgOpts.SetHeader("To", email)
		mutex.Lock()
		for {
			err := gomail.Send(smtpConn.Conn, msgOpts)
			if err == nil {
				pgBar.Add(1)
				smtpConn.NewConn = false
				break
			} else if smtpConn.NewConn && smtpConn.NumErrors >= 2 {
				color.New(color.FgRed).Printf("\nThe SMTP has been rate limited. Try again later!\n")
				os.Exit(1)
			} else if err != nil {
				if smtpConn.NewConn {
					smtpConn.NumErrors++
				}
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
