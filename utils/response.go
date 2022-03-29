package utils

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func GetSuccessResponse(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Content-Type":                     "application/json",
		},
	}
}

func GetErrorResponse(code int, message string, err error) events.APIGatewayProxyResponse {
	if err != nil {
		fmt.Println(err)
	}

	type Error struct {
		Message string `json:"message"`
	}

	type Body struct {
		Error Error `json:"error"`
	}

	error := Error{message}
	bytes, _ := json.Marshal(Body{error})
	return events.APIGatewayProxyResponse{
		Body:       string(bytes),
		StatusCode: code,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Content-Type":                     "application/json",
		},
	}
}
