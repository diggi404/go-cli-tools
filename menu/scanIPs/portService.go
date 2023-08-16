package menu

import (
	"errors"
	"net"
	"strings"
	"time"
)

// get the services running on each port connection created.
func GetServiceInfo(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return "Timeout", err
		} else {
			return "Timeout", err
		}
	}
	data := string(buffer[:n])
	return data, nil
}

// clear the port service info off any newline characters.
func SanitizeServiceInfo(serviceInfo string) string {
	sanitized := strings.ReplaceAll(serviceInfo, "\n", " ")
	sanitized = strings.TrimSpace(sanitized)
	return sanitized
}
