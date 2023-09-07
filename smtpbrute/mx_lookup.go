package smtpbrute

import (
	"errors"
	"net"
	"strings"
)

func LookupDomain(creds []string) ([]string, error) {
	username, password := creds[0], creds[1]
	domainSplit := strings.Split(username, "@")
	if len(domainSplit) != 2 {
		err := errors.New("invalid email format")
		return nil, err
	}
	domain := domainSplit[1]
	record, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}
	var host []string
	for _, v := range record {
		host = append(host, v.Host)
	}
	newHost := strings.TrimSuffix(host[0], ".")
	return []string{username, password, newHost}, nil
}
