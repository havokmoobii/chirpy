package auth

import(
	"time"
	"log"
	"net/http"
	"errors"

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
		log.Printf("%s", err)
		return uuid.Nil, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("%s", err)
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(subject)
	if err != nil {
		log.Printf("%s", err)
		return uuid.Nil, err
	}

	return userID, nil 
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer, ok := headers["Authorization"]
	if !ok {
		return "", errors.New("Missing Authorization header key")
	}
	malformedErr := errors.New("Malformed Authorization header value. Expecting: bearer <token>")
	if len(bearer) <= 1{
		return "", malformedErr
	}
	if bearer[0] != "bearer" {
		return "", malformedErr
	}
	return bearer[1], nil
}