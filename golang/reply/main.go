package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	telegram "example.com/telegram"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	HTTPMethodNotSupported = errors.New("HTTP method not supported")
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Body size = %d. \n", len(request.Body))
	logRequest := fmt.Sprintf("%+v", request)
	log.Printf("Request: %s", strings.ReplaceAll(logRequest, "\n", `\n`))

	s3BucketId := os.Getenv("s3_bucket_id")
	languageCode := os.Getenv("language")
	telegram.Reply(ctx, request.Body, s3BucketId, languageCode)
	return events.APIGatewayProxyResponse{Body: "POST", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
