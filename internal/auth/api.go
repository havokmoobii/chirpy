package auth

import(
	"net/http"
	"strings"
	"errors"
)

func GetAPIKey(headers http.Header) (string, error) {
	bearer, ok := headers["Authorization"]
	if !ok {
		return "", errors.New("Missing Authorization header key")
	}
	malformedErr := errors.New("Malformed Authorization header value. Expecting: ApiKey <token>")
	if len(bearer) == 0 {
		return "", malformedErr
	}
	split := strings.Split(bearer[0], " ")
	if len(split) != 2 {
	return "", malformedErr
	}
	if split[0] != "ApiKey" {
		return "", malformedErr
	}
	return split[1], nil
}