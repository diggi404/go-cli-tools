package webmail

import (
	"bytes"
	"fmt"
	"net/http"
)

func MakeRequest(creds []string) ([]string, error) {
	username, password, targetURL := creds[0], creds[1], creds[2]
	targetURL = fmt.Sprintf("%s/login/?login_only=1", targetURL)
	payloadStr := fmt.Sprintf("user=%s&pass=%s", username, password)
	payload := []byte(payloadStr)
	contentType := "application/x-www-form-urlencoded"
	res, err := http.Post(targetURL, contentType, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 200 {
		return creds, nil
	}
	return nil, err
}
