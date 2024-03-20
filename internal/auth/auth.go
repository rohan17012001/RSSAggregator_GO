package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey returns the API key from the headers
// Example:
// Authorization:ApiKey {insert API key here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no API key included in the header")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed auth header")
	}
	return vals[1], nil
}
