package main

import (
	"context"
	"log"
	"os"

	telegram "example.com/telegram"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, event MyEvent) error {
	s3BucketId := os.Getenv("s3_bucket_id")
	tokenParameterName := os.Getenv("token_parameter_name")
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	telegram.Send(ctx, cfg, s3BucketId, tokenParameterName)
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
