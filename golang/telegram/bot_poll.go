package telegram

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/aws/aws-sdk-go-v2/aws"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getSendPollConfig(dictionary []string, correctLineNumber int) (int, tgbotapi.SendPollConfig) {
	log.Printf("getSendPollConfig correctLineNumber %d\n", correctLineNumber)
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
	fileReader FileReader,
	languageCode string,
	chatID int64,
	tokenParameterName string,
) (*Chat, map[string]string, error) {
	batchId := -1
	correctWordId := -1
	correctLineNumber := -1
	thisChat, err := GetChat(chatID)
	if err == nil {
		if language, ok := thisChat.Languages[languageCode]; ok {
			if len(language.WrongAnswers) > 0 {
				correctWordId = extractWordId(language.WrongAnswers)
			} else if len(language.RightAnswers) > 0 && rand.Float32() > .5 {
				correctWordId = extractWordId(language.RightAnswers)
			}
		}
	}

	if correctWordId != -1 {
		log.Printf("correctWordId %d batchId = %d\n", correctWordId, batchId)
		batchId = int(correctWordId / 100)
		correctLineNumber = correctWordId % 100
	}
	dictionary, batchId, err := GetWords(fileReader, languageCode, batchId)
	if err != nil {
		log.Printf("Error while getting words: %s", err)
		return thisChat, nil, err
	}

	correctLineNumber, sendPollConfig := getSendPollConfig(dictionary, correctLineNumber)
	sendPollConfig.BaseChat = tgbotapi.BaseChat{
		ChatID: chatID,
	}

	params, err := Params(sendPollConfig)
	if err != nil {
		log.Printf("Error while getting params: %s", err)
		return thisChat, nil, err
	}

	params["method"] = "sendPoll"
	return thisChat, params, nil
}
