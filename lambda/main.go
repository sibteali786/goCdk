package main

import (
	"fmt"
	"lambda-func/app"
	"lambda-func/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"` // its called tagging
}

// Take in a payload and do something wit it
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username cannot be empty ")
	}

	return fmt.Sprintf("Successfully called by - %s\n", event.Username), nil
}

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "This is protected path",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	myApp := app.NewApp() // we will use it later
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return myApp.ApiHandler.RegisterUserHandler(request)
		case "/login":
			return myApp.ApiHandler.LoginUser(request)
		case "/protected":
			return middleware.ValidateJWTMiddleware(ProtectedHandler)(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not found ",
				StatusCode: 404,
			}, nil
		}
	})
}
