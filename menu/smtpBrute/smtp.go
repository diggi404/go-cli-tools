package smtpbrute

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

func ConnectSMTP(credentials []string, testEmail string) {
	username, password, host := credentials[0], credentials[1], credentials[2]
	dialer := gomail.NewDialer(host, 587, username, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", host)
	mailer.SetHeader("To", testEmail)
	mailer.SetHeader("Subject", "SMTP Test")
	mailer.SetBody("text/plain", "This is the body of the email.")
	err := dialer.DialAndSend(mailer)
	if err != nil {
		return
	}
	fmt.Println("email sent")
}
