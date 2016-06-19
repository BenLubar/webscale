// Package proxy contains functions useful for dealing with reverse proxies.
package proxy // import "github.com/BenLubar/webscale/proxy"

import (
	"net"
	"net/http"
	"strings"
)

var privateNets = func() []*net.IPNet {
	ranges := []string{
		"10.0.0.0/8",     // RFC 1918 private
		"127.0.0.0/8",    // RFC 990  loopback
		"169.254.0.0/16", // RFC 3927 link-local
		"172.16.0.0/12",  // RFC 1918 private
		"192.168.0.0/16", // RFC 1918 private
		"::1/128",        // RFC 4291 loopback
		"fc00::/7",       // RFC 4291 link-local
	}
	nets := make([]*net.IPNet, len(ranges))
	for i, cidr := range ranges {
		_, n, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}

		nets[i] = n
	}
	return nets
}()

// IsPrivateIP returns true if the IP is a loopback, link-local, or private
// network IP address.
func IsPrivateIP(ip net.IP) bool {
	for _, n := range privateNets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// RequestIP returns the IP address of the client that sent the request,
// assuming that private IP addresses are trusted reverse proxies that set the
// X-Forwarded-For header.
func RequestIP(r *http.Request) net.IP {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	ip := net.ParseIP(host)
	if ip != nil && !IsPrivateIP(ip) {
		return ip
	}

	candidates := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	for i := len(candidates) - 1; i >= 0; i-- {
		if candidate := net.ParseIP(strings.TrimSpace(candidates[i])); candidate != nil {
			ip = candidate
			if !IsPrivateIP(ip) {
				return ip
			}
		}
	}
	return ip
}
