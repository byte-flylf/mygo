package algs

import (
	"fmt"
	"strconv"
	"strings"
)

// Convert uint to net.IP
func inet_ntoa(ipnr uint) string {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return fmt.Sprintf("%d.%d.%d.%d", bytes[0], bytes[1], bytes[2], bytes[3])
}

// Convert net.IP to uint
func inet_aton(ipstr string) uint {
	bits := strings.Split(ipstr, ".")

	b0, _ := strconv.Atoi(bits[3])
	b1, _ := strconv.Atoi(bits[2])
	b2, _ := strconv.Atoi(bits[1])
	b3, _ := strconv.Atoi(bits[0])

	var sum uint
	sum += uint(b0) << 24
	sum += uint(b1) << 16
	sum += uint(b2) << 8
	sum += uint(b3)

	return sum
}
