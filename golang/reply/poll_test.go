package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	main "example.com/main"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestSerializer(t *testing.T) {

	tests := []struct {
		poll   *main.Poll
		expect string
	}{
		{
			poll: &main.Poll{
				ChatID:         123,
				WordLineNumber: 456,
				Language:       "sh",
				CreateEpoch:    789,
			},
			expect: "{\"ChatID\":123,\"WordLineNumber\":456,\"Language\":\"sh\",\"CreateEpoch\":789}",
		},
	}

	for _, test := range tests {
		data, err := json.Marshal(test.poll)
		assert.NoError(t, err)
		assert.Equal(t, test.expect, string(data))
	}

}

type MockS3Client struct {
}

func (s3Client *MockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return nil, nil
}

func TestPut(t *testing.T) {
	tests := []struct {
		poll   *main.Poll
		expect error
	}{
		{
			poll: &main.Poll{
				ChatID:         123,
				WordLineNumber: 456,
				Language:       "sh",
				CreateEpoch:    789,
			},
			expect: nil,
		},
		{
			poll: &main.Poll{
				WordLineNumber: 456,
				Language:       "sh",
				CreateEpoch:    789,
			},
			expect: fmt.Errorf("poll.ChatID is 0"),
		},
		{
			poll: &main.Poll{
				ChatID:      123,
				Language:    "sh",
				CreateEpoch: 789,
			},
			expect: fmt.Errorf("poll.WordLineNumber is 0"),
		},
		{
			poll: &main.Poll{
				ChatID:         123,
				WordLineNumber: 456,
				CreateEpoch:    789,
			},
			expect: fmt.Errorf("poll.Language is empty"),
		},
		{
			poll: &main.Poll{
				ChatID:         123,
				WordLineNumber: 456,
				Language:       "sh",
			},
			expect: fmt.Errorf("poll.CreateEpoch is 0"),
		},
	}

	for _, test := range tests {
		err := main.PutPoll(
			context.TODO(),
			&MockS3Client{},
			"",
			"",
			test.poll,
		)
		assert.Equal(t, test.expect, err)
	}

}
