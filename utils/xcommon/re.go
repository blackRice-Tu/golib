package xcommon

import (
	"regexp"
)

const (
	Ipv4Pattern = `^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`
)

func IsIpv4(ipv4 string) bool {
	regex := regexp.MustCompile(Ipv4Pattern)
	return regex.MatchString(ipv4)
}
