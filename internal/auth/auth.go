package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	id "github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "workoutBuilder-access"
)

func MakeJWT(userId id.UUID, tokenSecret string) (string, error) {

	signingKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour).UTC()),
		Subject:   userId.String(),
	})

	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (id.UUID, error) {

	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(key *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		fmt.Println(err)
		return id.Nil, err
	}

	userIdString, err := token.Claims.GetSubject()
	if err != nil {
		return id.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return id.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return id.Nil, errors.New("invalid Issuer")
	}

	uuID, err := id.Parse(userIdString)
	if err != nil {
		return id.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return uuID, nil
}

func GetBearerToken(headers http.Header) (string, error) {

	rawAuthString := headers.Get("Authorization")

	if rawAuthString == "" {
		return "", errors.New("authorization header was not provided")
	}

	splitAuthSlice := strings.Split(rawAuthString, "Bearer ")
	if len(splitAuthSlice) != 2 {
		return "", errors.New("badly formed authorization header")
	}

	trimmedToken := strings.TrimSpace(splitAuthSlice[1])

	return trimmedToken, nil

}
