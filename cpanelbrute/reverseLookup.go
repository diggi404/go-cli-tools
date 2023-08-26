package cpanelbrute

import "net"

func LookupIP(ip string) ([]string, error) {
	domains, err := net.LookupAddr(ip)
	if err != nil {
		return nil, err
	}
	return domains, nil
}
