package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	url := req.QueryStringParameters["url"]

	httpRes, err := http.Get(url)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Request error",
			StatusCode: 400,
		}, nil
	}

	defer httpRes.Body.Close()
	byteArray, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Request error",
			StatusCode: 400,
		}, nil
	}

	r1 := regexp.MustCompile("<title>(.+)<\\/title>")
	result := r1.FindString(string(byteArray))
	r2 := regexp.MustCompile("<title>")
	result = r2.ReplaceAllString(result, "")
	r3 := regexp.MustCompile("<\\/title>")
	result = r3.ReplaceAllString(result, "")

	return events.APIGatewayProxyResponse{
		Body:       "{\"title\": \"" + result + "\"}",
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Content-Type":                     "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
