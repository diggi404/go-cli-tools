package bomber

import "gopkg.in/gomail.v2"

func CreateSMTPConn(smtpCreds SmtpOpts) (gomail.SendCloser, error) {
	dialer := gomail.NewDialer(smtpCreds.Host, 587, smtpCreds.Username, smtpCreds.Password)
	conn, err := dialer.Dial()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
