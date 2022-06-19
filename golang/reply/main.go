package main

import (
	"context"
	"encoding/json"
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
	params, err := telegram.Reply(ctx, &OsFileReader{}, request.Body, languageCode, tokenParameterName)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
