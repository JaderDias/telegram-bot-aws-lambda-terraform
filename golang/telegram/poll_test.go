package telegram_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	telegram "example.com/telegram"

	"github.com/stretchr/testify/assert"
)

func TestPutPoll(t *testing.T) {
	clock := main.MockClock{}
	thirtyDays := clock.Now().Add(time.Hour * 24 * 30)
	tests := []struct {
		poll          *main.Poll
		expectErr     error
		expectBody    string
		expectExpires *time.Time
	}{
		{
			poll: &main.Poll{
				ChatID:   123,
				WordId:   456,
				Language: "sh",
			},
			expectErr:     nil,
			expectBody:    "{\"ChatID\":123,\"WordId\":456,\"Language\":\"sh\"}",
			expectExpires: &thirtyDays,
		},
		{
			poll: &main.Poll{
				WordId:   456,
				Language: "sh",
			},
			expectErr: fmt.Errorf("poll.ChatID is 0"),
		},
		{
			poll: &main.Poll{
				ChatID:   123,
				Language: "sh",
			},
			expectErr: fmt.Errorf("poll.WordId is 0"),
		},
		{
			poll: &main.Poll{
				ChatID: 123,
				WordId: 456,
			},
			expectErr: fmt.Errorf("poll.Language is empty"),
		},
	}

	for _, test := range tests {
		s3Client := main.MockS3Client{}
		err := main.PutPoll(
			context.TODO(),
			&s3Client,
			"",
			"",
			test.poll,
			&clock,
		)
		assert.Equal(t, test.expectErr, err)
		if test.expectErr == nil {
			assert.Equal(t, test.expectBody, main.ToString(s3Client.PutObjectInput.Body))
			assert.Equal(t, test.expectExpires, s3Client.PutObjectInput.Expires)
		}
	}

}

func TestGetPoll(t *testing.T) {
	tests := []struct {
		key        string
		expectPoll *main.Poll
		expectErr  error
		expectBody string
	}{
		{
			key:       "0",
			expectErr: fmt.Errorf("contentLength is 0"),
		},
		{
			key: "1",
			expectPoll: &main.Poll{
				ChatID:   123,
				WordId:   456,
				Language: "sh",
			},
			expectErr: nil,
		},
	}

	for _, test := range tests {
		s3Client := main.MockS3Client{}
		actual, err := main.GetPoll(
			context.TODO(),
			&s3Client,
			"",
			test.key,
		)
		assert.Equal(t, test.expectErr, err)
		if err == nil {
			assert.Equal(t, test.expectPoll, actual)
		}
	}

}
