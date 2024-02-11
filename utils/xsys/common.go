package xsys

import (
	"net"

	"github.com/pkg/errors"
)

// GetHostIpv4 ...
func GetHostIpv4() (string, error) {
	addressList, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.Wrap(err, "net.InterfaceAddrs")
	}
	for _, address := range addressList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("can not find host ip")
}
