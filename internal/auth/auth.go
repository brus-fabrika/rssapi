package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Get the API key from the request header
// Example:
// Authorization: ApiKey 1234567890

func GetApiKey(headers http.Header) (string, error) {

	// Get the API key from the request header
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is missing")
	}

	// Split the header into two parts: "ApiKey" and "1234567890"
	val := strings.SplitN(authHeader, " ", 2)
	if len(val) != 2 {
		return "", errors.New("Authorization header is invalid")
	}

	if val[0] != "ApiKey" {
		return "", errors.New("Authorization header is invalid")
	}

	return val[1], nil
}
