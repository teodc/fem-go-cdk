package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"lambda/app"
	"lambda/middleware"
	"net/http"
)

func ProtectedHandlerTest(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "{ \"ok\": true, \"message\": \"protected route\" }",
	}, nil
}

func main() {
	myApp := app.NewApp()
	handlerFunc := func(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		switch req.Path {
		case "/users/register":
			return myApp.ApiUserHandler.RegisterUser(req)
		case "/users/login":
			return myApp.ApiUserHandler.LoginUser(req)
		case "/users/protected":
			return middleware.ValidateJWT(ProtectedHandlerTest)(req)
		default:
			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       "{ \"ok\": false, \"message\": \" route not found\" }",
			}, nil
		}
	}
	lambda.Start(handlerFunc)
}
