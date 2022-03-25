package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func PutFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

func GetMd5(str string) string {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

type Image struct {
	Key        string `json:"key"`
	Type       string `json:"type"`
	Attachment string `json:"attachment"`
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucket := "melt-storage" // TODO: configから参照するようにする

	var image Image
	if err := json.Unmarshal([]byte(req.Body), &image); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Request Body parse error",
			StatusCode: 400,
		}, nil
	}

	// 先頭の ~;base64, まではファイルデータとして不要なので空文字で置換する
	r := regexp.MustCompile("^data:\\w+\\/\\w+;base64,")
	fileData := r.ReplaceAllString(image.Attachment, "")

	dec, err := base64.StdEncoding.DecodeString(fileData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Decode error",
			StatusCode: 400,
		}, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Configuration error",
			StatusCode: 500,
		}, nil
	}

	client := s3.NewFromConfig(cfg)

	extension := strings.Split(image.Key, ".")[1]
	key := "images/" + GetMd5(image.Key) + "." + extension

	ioReaderData := strings.NewReader(string(dec))
	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   ioReaderData,
	}

	_, err = PutFile(context.TODO(), client, input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Got error uploading file",
			StatusCode: 400,
		}, nil
	}

	url := "https://s3-ap-northeast-1.amazonaws.com/" + bucket + "/" + key
	return events.APIGatewayProxyResponse{
		Body:       "{\"name\": \"" + image.Key + "\", \"url\": \"" + url + "\"}",
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
