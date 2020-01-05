package networking

import (
	//"fmt"
	"net"
	//"strconv"
	"strings"
)

// GetCANInterfaces ...
func GetCANInterfaces() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var filtered []net.Interface
	for _, iface := range ifaces {
		if iface.Flags == 1 && (strings.HasPrefix(iface.Name, "can") || strings.HasPrefix(iface.Name, "vcan")) {
			filtered = append(filtered, iface)
		}
	}
	if len(filtered) == 0 {
		return nil, nil
	}
	return filtered, nil
}
