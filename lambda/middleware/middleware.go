package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	"lambda/auth"
)

func ValidateJWT(next func(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)) func(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return func(req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		token, err := extractTokenFromHeaders(&req.Headers)
		if err != nil {
			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "{ \"ok\": false, \"message\": \"cannot extract token from headers\" }",
			}, fmt.Errorf("error while extracting token from headers: %w", err)
		}

		jwtToken, err := auth.ParseJWTToken(token)
		if err != nil {
			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "{ \"ok\": false, \"message\": \"error while parsing JWT token\" }",
			}, fmt.Errorf("error while parsing JWT token: %w", err)
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "{ \"ok\": false, \"message\": \"invalid claims\" }",
			}, fmt.Errorf("invalid JWT claims")
		}

		expires := int64(claims["expires"].(float64))
		if expires < time.Now().Unix() {
			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "{ \"ok\": false, \"message\": \"token expired\" }",
			}, fmt.Errorf("JWT token has expired")
		}

		return next(req)
	}
}

func extractTokenFromHeaders(headers *map[string]string) (string, error) {
	if headers == nil {
		return "", fmt.Errorf("no headers provided")
	}

	authHeader, ok := (*headers)["Authorization"]
	if !ok {
		return "", fmt.Errorf("no Authorization header found")
	}

	bearerPrefix := "Bearer "
	splitAuthHeader := strings.Split(authHeader, bearerPrefix)
	if len(splitAuthHeader) != 2 {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	token := strings.TrimSpace(splitAuthHeader[1])
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}
