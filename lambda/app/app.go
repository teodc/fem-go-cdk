package app

import (
	"lambda/api"
	"lambda/database"
)

type App struct {
	ApiUserHandler *api.UserHandler
}

func NewApp() *App {
	store := database.NewDynamoDBStore()

	return &App{
		ApiUserHandler: api.NewUserHandler(store),
	}
}
