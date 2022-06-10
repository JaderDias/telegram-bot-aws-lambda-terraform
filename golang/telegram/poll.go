package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Poll struct {
	ChatID   int64
	WordId   int
	Language string
}

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (*RealClock) Now() time.Time { return time.Now() }

type MockClock struct{}

func (*MockClock) Now() time.Time { return time.Unix(1654822800, 0) }

func PutPoll(
	ctx context.Context,
	s3Client S3Client,
	s3BucketId string,
	pollId string,
	poll *Poll,
	clock Clock,
) error {
	if poll == nil {
		return fmt.Errorf("poll is nil")
	}
	if poll.ChatID == 0 {
		return fmt.Errorf("poll.ChatID is 0")
	}
	if poll.WordId == 0 {
		return fmt.Errorf("poll.WordId is 0")
	}
	if poll.Language == "" {
		return fmt.Errorf("poll.Language is empty")
	}
	key := fmt.Sprintf("poll/%s", pollId)
	data, err := json.Marshal(poll)
	if err != nil {
		return err
	}
	expires := clock.Now().Add(time.Hour * 24 * 30) // 30 days
	return PutObject(
		ctx,
		s3Client,
		s3BucketId,
		key,
		data,
		&expires,
	)
}

func GetPoll(
	ctx context.Context,
	s3Client S3Client,
	s3BucketId string,
	pollId string,
) (*Poll, error) {
	key := fmt.Sprintf("poll/%s", pollId)
	content, err := GetObject(ctx, s3Client, s3BucketId, key)
	if err != nil {
		return nil, err
	}

	var value Poll
	err = json.Unmarshal(content, &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
