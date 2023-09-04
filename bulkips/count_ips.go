package bulkips

import "net"

// CountTotalIPs Calculate the total number of IPs that can be generated within the given range;
// This is used or the cli progress bar.
func CountTotalIPs(startIP, endIP net.IP) int {
	start := IpToInt(startIP)
	end := IpToInt(endIP)
	return end - start + 1
}
