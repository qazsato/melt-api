package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	bucket := os.Getenv("S3_BUCKET_NAME")

	var image Image
	if err := json.Unmarshal([]byte(req.Body), &image); err != nil {
		return GetErrorResponse(400, "key, type, attachment are required"), nil
	}

	// 先頭の ~;base64, まではファイルデータとして不要なので空文字で置換する
	r := regexp.MustCompile("^data:\\w+\\/\\w+;base64,")
	fileData := r.ReplaceAllString(image.Attachment, "")
	dec, err := base64.StdEncoding.DecodeString(fileData)
	if err != nil {
		return GetErrorResponse(500, "Internal Server Error"), nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return GetErrorResponse(500, "Internal Server Error"), nil
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
		return GetErrorResponse(500, "Internal Server Error"), nil
	}

	url := "https://s3-ap-northeast-1.amazonaws.com/" + bucket + "/" + key
	body := map[string]string{
		"name": image.Key,
		"url":  url,
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
