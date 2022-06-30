package telegram

import (
	"context"
	"encoding/json"
	"fmt"
)

type Language struct {
	RightAnswers              map[int]bool
	WrongAnswers              map[int]bool
	SubscriptionPeriodSeconds int
	SubscriptionStartEpoch    int
	SubscriptionEndEpoch      int
}

type Chat struct {
	Languages map[string]Language
}

func PutChat(
	ctx context.Context,
	s3Client S3Client,
	s3BucketId string,
	chatId int64,
	chat *Chat,
) error {
	if chat == nil {
		return fmt.Errorf("chat is nil")
	}
	key := fmt.Sprintf("chat/%d", chatId)
	data, err := json.Marshal(chat)
	if err != nil {
		return err
	}
	return PutObject(
		ctx,
		s3Client,
		s3BucketId,
		key,
		data,
		nil,
	)
}

func GetChat(
	ctx context.Context,
	s3Client S3Client,
	s3BucketId string,
	chatId int64,
) (*Chat, error) {
	key := fmt.Sprintf("chat/%d", chatId)
	return GetChatFromKey(ctx, s3Client, s3BucketId, key)
}

func GetChatFromKey(
	ctx context.Context,
	s3Client S3Client,
	s3BucketId string,
	key string,
) (*Chat, error) {
	content, err := GetObject(ctx, s3Client, s3BucketId, key)
	if err != nil {
		return nil, err
	}

	var value Chat
	err = json.Unmarshal(content, &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
