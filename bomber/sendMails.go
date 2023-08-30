package bomber

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"gopkg.in/gomail.v2"
)

func SendMail(articleChunks <-chan Article, workingChan <-chan []int, wg *sync.WaitGroup, mutex *sync.Mutex, smtpConn *gomail.SendCloser, msgOptions *gomail.Message, smtpCreds *SmtpOpts, sucBar *progressbar.ProgressBar, failsBar *progressbar.ProgressBar) {
	defer wg.Done()
	article := <-articleChunks
	rounds := <-workingChan
	for range rounds {
		if article.Author == "" {
			article.Author = "Henry Sams"
		}
		rand.New(rand.NewSource(time.Now().UnixNano()))
		randomInt := rand.Intn(1000)
		randNum := strconv.Itoa(randomInt)
		msgOptions.SetAddressHeader("From", smtpCreds.Username, article.Author)
		msgOptions.SetHeader("Subject", article.Title+randNum)
		msgOptions.SetBody("text/plain", article.Description)
		mutex.Lock()
		err := gomail.Send(*smtpConn, msgOptions)
		if err == nil {
			sucBar.Add(1)
		} else {
			failsBar.Add(1)
		}
		mutex.Unlock()
	}
}
