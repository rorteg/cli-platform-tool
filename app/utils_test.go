package app

import (
	"net"
	"strings"
)

// Helper function to check if an IP is valid
func isValidIP(ip string) bool {
	trimmedIP := strings.TrimSpace(ip)
	return net.ParseIP(trimmedIP) != nil
}
