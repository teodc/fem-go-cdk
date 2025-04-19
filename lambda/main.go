package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"lambda/app"
	"net/http"
)

func main() {
	myApp := app.NewApp()
	handlerFunc := func(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		switch req.Path {
		case "/users/register":
			return myApp.ApiUserHandler.RegisterUser(req)
		case "/users/login":
			return myApp.ApiUserHandler.LoginUser(req)
		default:
			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       "Not Found",
			}, nil
		}
	}
	lambda.Start(handlerFunc)
}
