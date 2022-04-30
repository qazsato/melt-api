package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djimenez/iconv-go"
	"github.com/qazsato/melt-api/utils"
	"github.com/saintfish/chardet"
)

func GetTitle(body io.ReadCloser) (string, error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	// 文字コードの判定
	// cf. https://qiita.com/koki_develop/items/dab4bcbb1df1271a17b6
	detector := chardet.NewTextDetector()
	detectorResult, err := detector.DetectBest(bytes)
	if err != nil {
		return "", err
	}

	// title タグのテキスト取得
	reader := strings.NewReader(string(bytes))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}
	title := ""
	doc.Find("head title").Each(func(i int, s *goquery.Selection) {
		title = s.Text()
	})

	// 文字コードをUTF-8に変換
	// cf. https://github.com/djimenez/iconv-go#converting-string-values
	convertedTitle, err := iconv.ConvertString(title, detectorResult.Charset, "utf-8")
	if err != nil {
		return "", err
	}

	return convertedTitle, nil
}

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

	url := req.QueryStringParameters["url"]

	if url == "" {
		return utils.GetErrorResponse(400, "url is required", nil), nil
	}

	httpRes, err := http.Get(url)
	if err != nil {
		return utils.GetErrorResponse(500, "Internal Server Error", err), nil
	}

	defer httpRes.Body.Close()
	title, err := GetTitle(httpRes.Body)
	if err != nil {
		return utils.GetErrorResponse(500, "Internal Server Error", err), nil
	}

	body := map[string]string{
		"title": title,
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
