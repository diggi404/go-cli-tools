package cpanel

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func MakeRequest(creds []string) ([]string, error) {
	url, username, password := creds[0], creds[1], creds[2]
	targetURL := fmt.Sprintf("%s/login/?login_only=1", url)
	payloadStr := fmt.Sprintf("user=%s&pass=%s", username, password)
	payload := []byte(payloadStr)
	contentType := "application/x-www-form-urlencoded"
	res, err := http.Post(targetURL, contentType, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 200 {
		return creds, nil
	} else {
		err := errors.New("invalid credentials")
		return nil, err
	}
}
