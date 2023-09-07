package mailer

import (
	"fmt"
	"sync"

	"github.com/schollz/progressbar/v3"
	"gopkg.in/gomail.v2"
)

func SendMail(emailChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, smtpConn *gomail.SendCloser, msgOpts *gomail.Message, bar *progressbar.ProgressBar, results *MailOutResults) {
	defer wg.Done()
	emailList := <-emailChunks
	for _, email := range emailList {
		msgOpts.SetHeader("To", email)
		mutex.Lock()
		err := gomail.Send(*smtpConn, msgOpts)
		if err == nil {
			results.Success += 1
			bar.Add(1)
		} else {
			fmt.Printf("err: %v\n", err)
			results.Fails += 1
			bar.Add(1)
		}
		mutex.Unlock()

	}

}
