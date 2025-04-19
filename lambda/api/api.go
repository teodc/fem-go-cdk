package api

import (
	"fmt"
	"lambda/database"
	"lambda/types"
)

type Handler struct {
	dbClient *database.DynamoDBClient
}

func NewApiHandler(dbClient *database.DynamoDBClient) *Handler {
	return &Handler{
		dbClient: dbClient,
	}
}

func (h *Handler) RegisterUser(payload *types.RegisterUserPayload) error {
	if payload.Username == "" || payload.Password == "" {
		return fmt.Errorf("missing username or password")
	}

	userAlreadyExists, err := h.dbClient.DoesUserExist(payload.Username)
	if err != nil {
		return fmt.Errorf("error while checking user existence: %w", err)
	}
	if userAlreadyExists {
		return fmt.Errorf("user already exists")
	}

	err = h.dbClient.CreateUser(payload.Username, payload.Password)
	if err != nil {
		return fmt.Errorf("error while creating user: %w", err)
	}

	return nil
}
