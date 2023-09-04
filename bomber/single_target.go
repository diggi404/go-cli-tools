package bomber

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"gopkg.in/gomail.v2"
)

func SingleBomb(articleChunks <-chan Article, workingChan <-chan []int, wg *sync.WaitGroup, mutex *sync.Mutex, smtpConn *gomail.SendCloser, msgOptions *gomail.Message, smtpCreds *SmtpOpts, pgBar *progressbar.ProgressBar) {
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
		err := gomail.Send(*smtpConn, msgOptions)
		if err == nil {
			pgBar.Add(1)
		}
		mutex.Unlock()
	}
}
