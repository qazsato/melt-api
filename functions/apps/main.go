package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/qazsato/melt-api/utils"
	"github.com/tidwall/gjson"
)

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	apiKey := req.QueryStringParameters["api_key"]
	if apiKey == "" {
		return utils.GetErrorResponse(401, "api_key is required", nil), nil
	}
	ok, err := utils.IsExistKey(apiKey)
	if err != nil {
		return utils.GetErrorResponse(500, "Internal Server Error", err), nil
	}
	if *ok == false {
		return utils.GetErrorResponse(401, "Unauthorized", nil), nil
	}

	url := "https://api.github.com/repos/qazsato/melt/releases/latest"
	httpRes, err := http.Get(url)
	if err != nil {
		return utils.GetErrorResponse(500, "Internal Server Error", err), nil
	}

	defer httpRes.Body.Close()
	byteArray, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return utils.GetErrorResponse(500, "Internal Server Error", err), nil
	}

	body := map[string]string{
		"version": gjson.Get(string(byteArray), "tag_name").String(),
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return utils.GetErrorResponse(500, "Internal Server Error", err), nil
	}

	return utils.GetSuccessResponse(string(bytes)), nil
}

func main() {
	lambda.Start(Handler)
}
