package types

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func (payload *RegisterUserPayload) IsValid() error {
	if payload.Username == "" {
		return fmt.Errorf("missing username")
	}

	if payload.Password == "" {
		return fmt.Errorf("missing password")
	}

	return nil
}

func (payload *LoginUserPayload) IsValid() error {
	if payload.Username == "" {
		return fmt.Errorf("missing username")
	}

	if payload.Password == "" {
		return fmt.Errorf("missing password")
	}

	return nil
}

func NewUser(payload *RegisterUserPayload) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("error while hashing password: %w", err)
	}

	u := &User{
		Username:     payload.Username,
		PasswordHash: string(hashedPassword),
	}

	return u, nil
}

func ValidateUserPassword(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return false
	}

	return true
}
