package telegram_test

import (
	"context"
	"testing"
	"time"

	telegram "example.com/telegram"

	"github.com/stretchr/testify/assert"
)

func TestPutChat(t *testing.T) {
	tests := []struct {
		chatId        int64
		chat          *main.Chat
		expectErr     error
		expectBody    string
		expectExpires *time.Time
	}{
		{
			chatId: 123,
			chat: &main.Chat{
				Languages: map[string]main.Language{
					"sh": {
						RightAnswers: map[int]bool{1: true, 2: true, 3: true},
						WrongAnswers: map[int]bool{4: true, 5: true, 6: true},
					},
				},
			},
			expectErr:     nil,
			expectBody:    "{\"Languages\":{\"sh\":{\"RightAnswers\":{\"1\":true,\"2\":true,\"3\":true},\"WrongAnswers\":{\"4\":true,\"5\":true,\"6\":true}}}}",
			expectExpires: nil,
		},
	}

	for _, test := range tests {
		s3Client := main.MockS3Client{}
		err := main.PutChat(
			context.TODO(),
			&s3Client,
			"",
			test.chatId,
			test.chat,
		)
		assert.Equal(t, test.expectErr, err)
		if err == nil {
			assert.Equal(t, test.expectBody, main.ToString(s3Client.PutObjectInput.Body))
			assert.Equal(t, test.expectExpires, s3Client.PutObjectInput.Expires)
		}
	}
}

func TestGetChat(t *testing.T) {
	tests := []struct {
		key        int64
		expectChat *main.Chat
		expectErr  error
		expectBody string
	}{
		{
			key: 1,
			expectChat: &main.Chat{
				Languages: map[string]main.Language{
					"sh": {
						RightAnswers: map[int]bool{1: true, 2: true, 3: true},
						WrongAnswers: map[int]bool{4: true, 5: true, 6: true},
					},
				},
			},
			expectErr: nil,
		},
	}

	for _, test := range tests {
		s3Client := main.MockS3Client{}
		actual, err := main.GetChat(
			context.TODO(),
			&s3Client,
			"",
			test.key,
		)
		assert.Equal(t, test.expectErr, err)
		if test.expectErr == nil {
			assert.Equal(t, test.expectChat, actual)
		}
	}

}
