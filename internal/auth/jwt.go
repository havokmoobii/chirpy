package auth

import(
	"time"
	"log"
	"net/http"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
	return []byte(tokenSecret), nil})
	if err != nil {
		return uuid.Nil, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil 
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer, ok := headers["Authorization"]
	if !ok {
		return "", errors.New("Missing Authorization header key")
	}
	malformedErr := errors.New("Malformed Authorization header value. Expecting: Bearer <token>")
	if len(bearer) == 0 {
		return "", malformedErr
	}
	split := strings.Split(bearer[0], " ")
	if len(split) != 2 {
	return "", malformedErr
	}
	if split[0] != "Bearer" {
		return "", malformedErr
	}
	return split[1], nil
}