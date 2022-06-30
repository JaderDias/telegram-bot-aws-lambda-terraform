package main

import (
	"context"
	"log"
	"os"

	telegram "example.com/telegram"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	if len(os.Args) < 4 {
		log.Println("Usage: ./main <tokenParameterName> <s3BucketId> <region>")
		os.Exit(1)
	}

	tokenParameterName := os.Args[1]
	s3BucketId := os.Args[2]
	awsRegion := os.Args[3]
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	telegram.Send(ctx, cfg, s3BucketId, tokenParameterName)
}
