package bulkips

import "net"

// Compares the starting and ending IPs. Used as a condition in the for loop generating the IPs.
func BytesCompare(a, b net.IP) int {
	for i := range a {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
