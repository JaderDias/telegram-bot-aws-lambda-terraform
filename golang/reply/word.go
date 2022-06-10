package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetWords(
	ctx context.Context,
	s3Client *s3.Client,
	s3BucketId string,
	language string,
) ([]string, int, error) {
	batch := rand.Intn(260)
	key := fmt.Sprintf("language/%s/%d", language, batch)
	bbuffer, err := GetObject(ctx, s3Client, s3BucketId, key)
	if err != nil {
		return nil, 0, err
	}

	return strings.Split(bbuffer.String(), "\n"), batch, nil
}
