package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	HTTPMethodNotSupported = errors.New("HTTP method not supported")
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: "POST", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
