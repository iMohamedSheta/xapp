package utils

import (
	"fmt"
	"net"
	"strings"
)

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}

		ipv4 := ipNet.IP.To4()
		if ipv4 == nil {
			continue
		}

		// return first non-loopback IPv4 (LAN IP)
		return ipv4.String(), nil
	}

	return "", fmt.Errorf("no LAN IP found")
}

// OriginHost extracts just the hostname from an Origin value like
// "http://192.168.1.42:5173".
func OriginHost(origin string) string {
	raw := origin
	if idx := strings.Index(raw, "://"); idx != -1 {
		raw = raw[idx+3:]
	}
	host, _, err := net.SplitHostPort(raw)
	if err != nil {
		return raw // no port
	}
	return host
}

// IsPrivateOrLoopback returns true when ip is an RFC-1918 private address or
// a loopback address (127.x.x.x / ::1).
func IsPrivateOrLoopback(input string) bool {
	ip := net.ParseIP(input)

	// If not an IP → try resolving (handles "localhost", domains)
	if ip == nil {
		ips, err := net.LookupIP(input)
		if err != nil || len(ips) == 0 {
			return false
		}
		ip = ips[0]
	}

	return ip.IsLoopback() || ip.IsPrivate()
}
