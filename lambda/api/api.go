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
	err := payload.IsValid()
	if err != nil {
		return fmt.Errorf("user registration payload not valid: %w", err)
	}

	exists, err := handler.store.DoesUserExist(payload.Username)
	if err != nil {
		return fmt.Errorf("error while checking user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("user already exists")
	}

	user, err := types.NewUser(payload)
	if err != nil {
		return err
	}

	err = handler.store.PersistUser(user)
	if err != nil {
		return fmt.Errorf("error while creating user: %w", err)
	}

	return nil
}
