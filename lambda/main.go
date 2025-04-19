package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"lambda/app"
)

func main() {
	newApp := app.NewApp()
	handler := newApp.ApiHandler.RegisterUser
	lambda.Start(handler)
}
