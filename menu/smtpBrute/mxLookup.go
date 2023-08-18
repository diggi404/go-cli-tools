package smtpbrute

import (
	"net"
	"strings"
)

func LookupDomain(creds []string) ([]string, error) {
	username, password := creds[0], creds[1]
	domainSplit := strings.Split(username, "@")
	domain := domainSplit[1]
	record, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}
	var host string
	for _, v := range record {
		host = v.Host
	}
	return []string{username, password, host}, nil
}
