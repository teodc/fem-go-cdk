package app

import (
	"lambda/api"
	"lambda/database"
)

type App struct {
	ApiHandler *api.Handler
}

func NewApp() *App {
	return &App{
		ApiHandler: api.NewApiHandler(
			database.NewDynamoDBClient(),
		),
	}
}
