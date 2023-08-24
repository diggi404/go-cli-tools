package bulkips

import (
	"net"
	"strconv"
	"strings"
)

// IpToInt This will convert an IP address[byte] into it's full integer value;
// Used in the CountTotalIPs function to calc the number of IPs within a range.
func IpToInt(ip net.IP) int {
	octets := strings.Split(ip.String(), ".")
	if len(octets) != 4 {
		return 0
	}

	var result int
	for _, octetStr := range octets {
		octet, err := strconv.Atoi(octetStr)
		if err != nil {
			return 0
		}
		result = result*256 + octet
	}

	return result
}
