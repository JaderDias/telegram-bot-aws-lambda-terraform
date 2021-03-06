package telegram

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getSendPollConfig(dictionary []string, correctLineNumber int) (int, tgbotapi.SendPollConfig) {
	options := [3]Word{}
	grammaticalClass := ""
	for i := 0; i < 3; {
		lineNumber := rand.Intn(len(dictionary))
		if i == 0 {
			if correctLineNumber != -1 {
				lineNumber = correctLineNumber
			} else {
				correctLineNumber = lineNumber
			}
		}
		options[i] = Parse(dictionary[lineNumber])
		if options[i].err != nil {
			log.Printf("Error while parsing line %d: %s", lineNumber, options[i].err)
			continue
		}
		if grammaticalClass == "" {
			grammaticalClass = options[i].grammaticalClass
		} else if options[i].grammaticalClass != grammaticalClass {
			continue
		}

		i++
	}

	correctAnswerIndex := rand.Intn(3)
	if correctAnswerIndex != 0 {
		aux := options[0]
		options[0] = options[correctAnswerIndex]
		options[correctAnswerIndex] = aux
	}

	correctAnswer := options[correctAnswerIndex]
	return correctLineNumber, tgbotapi.SendPollConfig{
		Type:     "quiz",
		Question: fmt.Sprintf("%s (%s)", correctAnswer.mainDefinition, correctAnswer.grammaticalClass),
		Options: []string{
			options[0].title,
			options[1].title,
			options[2].title,
		},
		CorrectOptionID: int64(correctAnswerIndex),
		IsAnonymous:     true,
	}
}

func extractWordId(answers map[int]bool) int {
	var id int
	for id, _ = range answers {
		break
	}

	delete(answers, id)
	return id
}

func BotSendPoll(
	ctx context.Context,
	cfg aws.Config,
	s3Client *s3.Client,
	s3BucketId string,
	languageCode string,
	chatID int64,
	tokenParameterName string,
) (*Chat, error) {
	correctWordId := -1
	thisChat, err := GetChat(ctx, s3Client, s3BucketId, chatID)
	if err == nil {
		if language, ok := thisChat.Languages[languageCode]; ok {
			correctWordId = GetCorrectWordId(language)
		}
	}
	err = BotSendPollWithCorrectWordId(ctx, cfg, s3Client, s3BucketId, languageCode, chatID, tokenParameterName, correctWordId)
	return thisChat, err
}

func GetCorrectWordId(
	language Language,
) int {
	if len(language.WrongAnswers) > 0 {
		return extractWordId(language.WrongAnswers)
	}

	if len(language.RightAnswers) > 0 && rand.Float32() > .5 {
		return extractWordId(language.RightAnswers)
	}

	return -1
}

func BotSendPollWithCorrectWordId(
	ctx context.Context,
	cfg aws.Config,
	s3Client *s3.Client,
	s3BucketId string,
	languageCode string,
	chatID int64,
	tokenParameterName string,
	correctWordId int,
) error {
	batchId := -1
	correctLineNumber := -1
	if correctWordId != -1 {
		log.Printf("correctWordId %d batchId = %d\n", correctWordId, batchId)
		batchId = int(correctWordId / 100)
		correctLineNumber = correctWordId % 100
	}
	dictionary, batchId, err := GetWords(ctx, s3Client, s3BucketId, languageCode, batchId)
	if err != nil {
		log.Printf("Error while getting words: %s", err)
		return err
	}

	correctLineNumber, sendPollConfig := getSendPollConfig(dictionary, correctLineNumber)
	sendPollConfig.BaseChat = tgbotapi.BaseChat{
		ChatID: chatID,
	}

	telegramBotTokens, err := GetTokens(ctx, cfg, tokenParameterName)
	if err != nil {
		log.Printf("unable to get telegram bot tokens, %v", err)
		return err
	}

	bot, err := GetBot(telegramBotTokens, languageCode)
	if err != nil {
		log.Printf("Error while creating bot: %s", err)
		return err
	}

	poll, err := bot.Send(sendPollConfig)
	if err != nil {
		log.Printf("Error while sending poll: %s", err)
		return err
	}

	err = PutPoll(
		ctx,
		s3Client,
		s3BucketId,
		poll.Poll.ID,
		&Poll{
			ChatID:   chatID,
			WordId:   (batchId * 100) + correctLineNumber,
			Language: languageCode,
		},
		&RealClock{},
	)
	if err != nil {
		log.Printf("Error while saving poll: %s", err)
		return err
	}

	return nil
}
