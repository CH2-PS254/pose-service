package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func ExtractBearerToken(header string) (string, error) {
	const BEARER_SCHEMA = "Bearer "

	if len(header) > len(BEARER_SCHEMA) && header[:len(BEARER_SCHEMA)] == BEARER_SCHEMA {
		return header[len(BEARER_SCHEMA):], nil
	}

	return "", errors.New("invalid authorization header")
}

func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return jwtKey, nil
	})

	return token, err
}
