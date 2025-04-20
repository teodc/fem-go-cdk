package api

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"lambda/database"
	"lambda/types"
	"net/http"
)

type UserHandler struct {
	store database.UserStore
}

func NewUserHandler(store database.UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (handler *UserHandler) RegisterUser(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var payload types.RegisterUserPayload

	err := json.Unmarshal([]byte(req.Body), &payload)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "{ \"ok\": false, \"message\": \"error while parsing request body\" }",
		}, fmt.Errorf("error while parsing request body: %w", err)
	}

	err = payload.IsValid()
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Body:       "{ \"ok\": false, \"message\": \"error while validating request\" }",
		}, fmt.Errorf("error while validating request: %w", err)
	}

	exists, err := handler.store.DoesUserExist(payload.Username)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "{ \"ok\": false, \"message\": \"error while checking if user exists\" }",
		}, fmt.Errorf("error while checking if user exists: %w", err)
	}
	if exists {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusConflict,
			Body:       "{ \"ok\": false, \"message\": \"user already exists\" }",
		}, nil
	}

	user, err := types.NewUser(&payload)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "{ \"ok\": false, \"message\": \"error while creating user\" }",
		}, fmt.Errorf("error while creating user: %w", err)
	}

	err = handler.store.PersistUser(user)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "{ \"ok\": false, \"message\": \"error while persisting user\" }",
		}, fmt.Errorf("error while persisting user: %w", err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "{ \"ok\": true, \"message\": \"user registered\" }",
	}, nil
}

func (handler *UserHandler) LoginUser(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var payload types.LoginUserPayload

	err := json.Unmarshal([]byte(req.Body), &payload)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "{ \"ok\": false, \"message\": \"error while parsing request body\" }",
		}, fmt.Errorf("error while parsing request body: %w", err)
	}

	err = payload.IsValid()
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Body:       "{ \"ok\": false, \"message\": \"error while validating request\" }",
		}, fmt.Errorf("error while validating request: %w", err)
	}

	user, err := handler.store.GetUser(payload.Username)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "{ \"ok\": false, \"message\": \"error while getting user\" }",
		}, fmt.Errorf("error while getting user: %w", err)
	}
	if user == nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "{ \"ok\": false, \"message\": \"user not found\" }",
		}, nil
	}

	passwordMatches := types.ValidateUserPassword(payload.Password, user.PasswordHash)
	if !passwordMatches {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "{ \"ok\": false, \"message\": \"invalid credentials\" }",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "{ \"ok\": true, \"message\": \"user logged in\" }",
	}, nil
}
