package main

import (
	"context"
	"errors"
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

type OsFileReader struct {
}

func (r *OsFileReader) ReadFile(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Body size = %d. \n", len(request.Body))
	logRequest := fmt.Sprintf("%+v", request)
	log.Printf("Request: %s", strings.ReplaceAll(logRequest, "\n", `\n`))

	languageCode := os.Getenv("language")
	tokenParameterName := os.Getenv("token_parameter_name")
	telegram.Reply(ctx, &OsFileReader{}, request.Body, languageCode, tokenParameterName)
	return events.APIGatewayProxyResponse{Body: "POST", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
