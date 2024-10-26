package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	subjectID, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}
	id, err := uuid.Parse(subjectID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed auth header")
	}
	return strings.TrimSpace(splitAuth[1]), nil
}
