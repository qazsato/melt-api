package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "{\"message\": \"GET /sites/metadata\"}",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
