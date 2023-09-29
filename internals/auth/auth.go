package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")
var ErrMalformedAuthorizationHeader = errors.New("malformed authorization header")

func GetBearerToken(headers http.Header, prefix string) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != prefix {
		return "", ErrMalformedAuthorizationHeader
	}
	return splitAuth[1], nil
}
