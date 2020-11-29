package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"

	uuid "github.com/google/uuid"
)

// GenerateID ..
func GenerateID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return "error_id"
	}
	return id.String()
}

// GetMD5Hash ..
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// SourceIP ..
func SourceIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("are you connected to the network?")
}
