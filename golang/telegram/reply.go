package telegram

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Reply(ctx context.Context, requestBody string) {
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

	ssmClient := ssm.NewFromConfig(cfg)

	parameterName := "telegram_bot_token"
	parameterOutput, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &parameterName,
		WithDecryption: true,
	})
	if err != nil {
		log.Printf("unable to get telegram bot token, %v", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(*parameterOutput.Parameter.Value)
	if err != nil {
		log.Printf("Error while creating bot: %s", err)
		return
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Debug = true

	s3BucketId := os.Getenv("s3_bucket_id")
	log.Printf("s3BucketId = %s\n", s3BucketId)
	s3Client := s3.NewFromConfig(cfg)

	if update.Message != nil { // If we got a message
		thisChat, err := BotSendPoll(
			ctx,
			s3Client,
			s3BucketId,
			bot,
			update.Message.Chat.ID,
		)
		if err != nil {
			log.Printf("Error while sending poll: %s", err)
			return
		}

		if thisChat == nil {
			PutChat(ctx, s3Client, s3BucketId, update.Message.Chat.ID, thisChat)
		}

	} else if update.Poll != nil {
		poll, err := GetPoll(ctx, s3Client, s3BucketId, update.Poll.ID)
		if err != nil {
			log.Printf("Error while retrieving poll: %s", err)
			return
		}

		thisChat, err := BotSendPoll(
			ctx,
			s3Client,
			s3BucketId,
			bot,
			poll.ChatID)
		if err != nil {
			log.Printf("Error while sending poll: %s", err)
			return
		}

		if thisChat == nil {
			thisChat = &Chat{}
		}
		language, ok := thisChat.Languages["sh"]
		if !ok {
			language = Language{}
			thisChat.Languages["sh"] = language
		}
		if update.Poll.Options[update.Poll.CorrectOptionID].VoterCount == 0 {
			language.WrongAnswers[poll.WordId] = true
		} else {
			language.RightAnswers[poll.WordId] = true
		}
		PutChat(ctx, s3Client, s3BucketId, poll.ChatID, thisChat)
	}
}
