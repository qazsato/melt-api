package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/qazsato/melt-api/utils"
)

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
	apiKey := req.QueryStringParameters["api_key"]
	if apiKey == "" {
		return utils.GetErrorResponse(401, "api_key is required"), nil
	}
	ok, err := utils.IsExistKey(apiKey)
	if err != nil {
		return utils.GetErrorResponse(500, "Internal Server Error"), nil
	}
	if *ok == false {
		return utils.GetErrorResponse(401, "Unauthorized"), nil
	}

	url := req.QueryStringParameters["url"]

	if url == "" {
		return utils.GetErrorResponse(400, "url is required"), nil
	}

	httpRes, err := http.Get(url)
	if err != nil {
		return utils.GetErrorResponse(500, "Internal Server Error"), nil
	}

	defer httpRes.Body.Close()
	byteArray, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return utils.GetErrorResponse(500, "Internal Server Error"), nil
	}

	body := map[string]string{
		"title": GetTitle(string(byteArray)),
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return utils.GetErrorResponse(500, "Internal Server Error"), nil
	}

	return utils.GetSuccessResponse(string(bytes)), nil
}

func main() {
	lambda.Start(Handler)
}
