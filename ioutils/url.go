package ioutils

import (
	"net/url"
	"strings"
)

// IsSecureScheme checks if the URL has a secure scheme.
func IsSecureScheme(u *url.URL) bool {
	secureSchemes := []string{"https", "tls", "mqtts", "wss"}
	for _, scheme := range secureSchemes {
		if strings.EqualFold(u.Scheme, scheme) {
			return true
		}
	}
	return false
}

// IsSecureUrl checks if the given URL is secure.
//
// Takes a string parameter `urlStr`.
// Returns a boolean indicating if the URL is secure and an error if any.
func IsSecureUrl(urlStr string) (bool, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false, err
	}
	return IsSecureScheme(u), err
}
