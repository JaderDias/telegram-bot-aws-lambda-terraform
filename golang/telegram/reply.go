package telegram

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Reply(ctx context.Context, fileReader FileReader, requestBody, languageCode string, tokenParameterName string) {
	var update tgbotapi.Update
	err := json.Unmarshal([]byte(requestBody), &update)
	if err != nil {
		log.Println(err)
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return
	}

	if update.Message != nil { // If we got a message
		thisChat, err := BotSendPoll(
			ctx,
			cfg,
			fileReader,
			languageCode,
			update.Message.Chat.ID,
			tokenParameterName,
		)
		if err != nil {
			log.Printf("Error while sending poll: %s", err)
			return
		}

		if thisChat == nil {
			PutChat(update.Message.Chat.ID, thisChat)
		}

	} else if update.Poll != nil {
		poll, err := GetPoll(update.Poll.ID)
		if err != nil {
			log.Printf("Error while retrieving poll: %s", err)
			return
		}

		thisChat, err := BotSendPoll(
			ctx,
			cfg,
			fileReader,
			languageCode,
			poll.ChatID,
			tokenParameterName,
		)
		if err != nil {
			log.Printf("Error while sending poll: %s", err)
			return
		}

		if thisChat == nil {
			thisChat = &Chat{}
		}
		if thisChat.Languages == nil {
			thisChat.Languages = make(map[string]Language)
		}
		language, ok := thisChat.Languages[languageCode]
		if !ok {
			language = Language{}
			thisChat.Languages[languageCode] = language
		}
		if update.Poll.Options[update.Poll.CorrectOptionID].VoterCount == 0 {
			if language.WrongAnswers == nil {
				language.WrongAnswers = make(map[int]bool)
			}
			language.WrongAnswers[poll.WordId] = true
		} else {
			if language.RightAnswers == nil {
				language.RightAnswers = make(map[int]bool)
			}
			language.RightAnswers[poll.WordId] = true
		}
		PutChat(poll.ChatID, thisChat)
	}
}
