package bomber

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/gomail.v2"
)

func SingleBomb(articleChunks <-chan Article, workingChan <-chan []int, wg *sync.WaitGroup, mutex *sync.Mutex, msgOptions *gomail.Message, smtpCreds SmtpOpts, pgBar *progressbar.ProgressBar, smtpConnIndex *int) {
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
		for i, conn := range smtpCreds.Conns {
			if i == *smtpConnIndex {
				err := gomail.Send(conn, msgOptions)
				if err == nil {
					pgBar.Add(1)
					break
				} else {
					color.New(color.FgRed).Printf("\nerr: %v\n", err)
					*smtpConnIndex++
					color.New(color.FgGreen).Printf("\nusing a different Connection\n")
				}
			}
		}
		mutex.Unlock()
	}
}
