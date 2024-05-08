package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

// gets the API Key
func GetApiKey(header http.Header) (string, error) {
	// gets the authorization
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	// splits the header to strip off the 'Bearer'
	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) < 2 || splitAuthHeader[0] != "ApiKey" {
		return "", errors.New("header not formatted correctly")
	}
	return splitAuthHeader[1], nil
}
