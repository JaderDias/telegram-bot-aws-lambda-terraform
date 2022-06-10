package telegram

import (
	"encoding/json"
	"fmt"
	"os"
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
	chatId int64,
	chat *Chat,
) error {
	if chat == nil {
		return fmt.Errorf("chat is nil")
	}
	key := fmt.Sprintf("/mnt/efs/chat/%d", chatId)
	data, err := json.Marshal(chat)
	if err != nil {
		return err
	}
	return os.WriteFile(key, data, 0644)
}

func GetChat(
	chatId int64,
) (*Chat, error) {
	key := fmt.Sprintf("/mnt/efs/chat/%d", chatId)
	content, err := os.ReadFile(key)
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
