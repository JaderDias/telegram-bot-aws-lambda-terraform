package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Poll struct {
	ChatID         int64
	WordLineNumber int
	Language       string
	CreateEpoch    int64
}

func PutPoll(
	ctx context.Context,
	s3Client S3Client,
	s3BucketId string,
	pollId string,
	poll *Poll,
) error {
	if poll == nil {
		return fmt.Errorf("poll is nil")
	}
	if poll.ChatID == 0 {
		return fmt.Errorf("poll.ChatID is 0")
	}
	if poll.WordLineNumber == 0 {
		return fmt.Errorf("poll.WordLineNumber is 0")
	}
	if poll.Language == "" {
		return fmt.Errorf("poll.Language is empty")
	}
	if poll.CreateEpoch == 0 {
		return fmt.Errorf("poll.CreateEpoch is 0")
	}
	key := fmt.Sprintf("poll/%s", pollId)
	data, err := json.Marshal(poll)
	if err != nil {
		return err
	}
	return PutObject(ctx, s3Client, s3BucketId, key, data)
}

func GetPoll(
	ctx context.Context,
	s3Client *s3.Client,
	s3BucketId string,
	pollId string,
) (*Poll, error) {
	key := fmt.Sprintf("poll/%s", pollId)
	bbuffer, err := GetObject(ctx, s3Client, s3BucketId, key)
	if err != nil {
		return nil, err
	}

	var value Poll
	err = json.Unmarshal(bbuffer.Bytes(), &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
