// Package validation implements networking-related validation-interfaces for the clt-package
package validation

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ValidateCANInterface validates the CAN interface
func ValidateCANInterface(s string) (bool, error) {
	if !(strings.HasPrefix(s, "can") || strings.HasPrefix(s, "vcan")) {
		return false, fmt.Errorf("'%s' is not a valid can interface", s)
	}
	_, err := net.InterfaceByName(s)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ValidateHost validates a given host
func ValidateHost(s string) (bool, error) {
	host := net.ParseIP(s)
	if host != nil {
		return true, nil
	}
	hostname, _ := net.LookupHost(s)
	if len(hostname) > 0 {
		return true, nil
	}
	return false, fmt.Errorf("'%s' does not seem to be a valid IP or Hostname", s)
}

// ValidatePort validates a given port
func ValidatePort(s string) (bool, error) {
	port, err := strconv.ParseInt(s, 10, 64)
	if err != nil || port < 0 || port > 65535 {
		return false, fmt.Errorf("'%s' is not a valid Port", s)
	}
	return true, nil
}
