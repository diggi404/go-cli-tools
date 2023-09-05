package bomber

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	"gopkg.in/gomail.v2"
)

func MultiTargetBomb(articleChunks <-chan Article, emailChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, msgOptions *gomail.Message, smtpCreds SmtpOpts, numBombs int, smtpConnIndex *int) {
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
}
