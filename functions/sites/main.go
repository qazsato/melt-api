package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func GetErrorResponse(code int, message string) events.APIGatewayProxyResponse {
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
	}
}

func GetTitle(html string) string {
	r1 := regexp.MustCompile("<title>(.+)<\\/title>")
	result := r1.FindString(html)
	r2 := regexp.MustCompile("<title>")
	result = r2.ReplaceAllString(result, "")
	r3 := regexp.MustCompile("<\\/title>")
	result = r3.ReplaceAllString(result, "")
	return result
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	url := req.QueryStringParameters["url"]

	if url == "" {
		return GetErrorResponse(400, "url is required"), nil
	}

	httpRes, err := http.Get(url)
	if err != nil {
		return GetErrorResponse(500, "Internal Server Error"), nil
	}

	defer httpRes.Body.Close()
	byteArray, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return GetErrorResponse(500, "Internal Server Error"), nil
	}

	body := map[string]string{
		"title": GetTitle(string(byteArray)),
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return GetErrorResponse(500, "Internal Server Error"), nil
	}

	return GetSuccessResponse(string(bytes)), nil
}

func main() {
	lambda.Start(Handler)
}
