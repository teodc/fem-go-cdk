package api

import (
	"fmt"
	"lambda/database"
	"lambda/types"
)

type UserHandler struct {
	store database.UserStore
}

func NewUserHandler(store database.UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (handler *UserHandler) RegisterUser(payload *types.RegisterUserPayload) error {
	if payload.Username == "" || payload.Password == "" {
		return fmt.Errorf("missing username or password")
	}

	userAlreadyExists, err := handler.store.DoesUserExist(payload.Username)
	if err != nil {
		return fmt.Errorf("error while checking user existence: %w", err)
	}
	if userAlreadyExists {
		return fmt.Errorf("user already exists")
	}

	err = handler.store.CreateUser(payload.Username, payload.Password)
	if err != nil {
		return fmt.Errorf("error while creating user: %w", err)
	}

	return nil
}
