package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"lambda/types"
	"time"
)

const (
	secret = "th3s3cr3t" // TODO: get from env
)

func MakeJWTToken(user *types.User) (string, error) {
	claims := jwt.MapClaims{
		"user":    user.Username,
		"expires": time.Now().Add(time.Hour * 1).Unix(), // valid for 1 hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error while signing token: %w", err)
	}

	return tokenString, nil
}

func ParseJWTToken(token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Here you can validate a bunch of things like the signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while parsing token: %w", err)
	}

	if !jwtToken.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return jwtToken, nil
}
