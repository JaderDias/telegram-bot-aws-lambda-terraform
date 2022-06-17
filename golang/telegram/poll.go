package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	pollId string,
	poll *Poll,
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
	key := fmt.Sprintf("/mnt/efs/poll/%s", pollId)
	data, err := json.Marshal(poll)
	if err != nil {
		return err
	}
	return os.WriteFile(key, data, 0644)
}

func GetPoll(
	pollId string,
) (*Poll, error) {
	log.Printf("GetPoll pollId = %s\n", pollId)
	key := fmt.Sprintf("/mnt/efs/poll/%s", pollId)
	content, err := os.ReadFile(key)
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
