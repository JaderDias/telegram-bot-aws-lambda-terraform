package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	HTTPMethodNotSupported = errors.New("HTTP method not supported")
)

func loadDictionary() ([]string, error) {
	file, err := os.Open("sh.csv")
	if err != nil {
		return nil, fmt.Errorf("Error while opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error while reading file: %s", err)
	}
	return words, nil
}

var titleMatcher = regexp.MustCompile(`^([^=]+)(=.*)$`)
var undesiredSections = regexp.MustCompile(`(?s)====?(?:Conjugation|Declension|Derived terms|Pronunciation)====?[^=]*`)
var mainDefinitionSearcher = regexp.MustCompile(`(?s)===([^=]+)===[^#]*# ([^\n]*)`)
var removeTransitiveness = regexp.MustCompile(`{{indtr\|[^}|]*\|([^}])}}\s*`)
var removeCurlyLink = regexp.MustCompile(`{{[^}]*[|=]([^|}=]+)}}`)
var removeSquareLink = regexp.MustCompile(`\[\[(?:[^|]*\|)?([^|\]]*)\]\]`)

type Word struct {
	title            string
	grammaticalClass string
	mainDefinition   string
	err              error
}

func Parse(s string) Word {
	titleSize := strings.Index(s, "=")
	if titleSize == 0 {
		return Word{
			err: errors.New("Invalid title"),
		}
	}

	title := s[:titleSize]

	// replace escaped line brakes with newline
	s = strings.Replace(s[titleSize:], "\\n", "\n", -1)

	// remove undesired sections
	s = undesiredSections.ReplaceAllString(s, "")

	section := mainDefinitionSearcher.FindStringSubmatch(s)
	if len(section) < 3 {
		return Word{
			err: errors.New("No mainDefinition found"),
		}
	}

	mainDefinition := removeTransitiveness.ReplaceAllString(section[2], "")
	mainDefinition = removeCurlyLink.ReplaceAllString(mainDefinition, "$1")
	mainDefinition = removeSquareLink.ReplaceAllString(mainDefinition, "$1")
	return Word{
		title:            title,
		grammaticalClass: section[1],
		mainDefinition:   mainDefinition,
	}
}

func getPoll(dictionary []string, correctLineNumber int) (int, tgbotapi.SendPollConfig) {
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

type poll struct {
	chatID         int64
	wordLineNumber int
}

type chat struct {
	wrongAnswers []int
	rightAnswers []int
}

func sendPoll(
	ctx context.Context,
	s3Client *s3.Client,
	s3BucketId string,
	dictionary []string,
	bot *tgbotapi.BotAPI,
	chatID int64,
) {
	correctLineNumber := -1
	correctLineNumber, sendPollConfig := getPoll(dictionary, correctLineNumber)
	sendPollConfig.BaseChat = tgbotapi.BaseChat{
		ChatID: chatID,
	}

	poll, err := bot.Send(sendPollConfig)
	if err != nil {
		log.Printf("Error while sending poll: %s", err)
		return
	}

	err = PutPoll(
		ctx,
		s3Client,
		s3BucketId,
		poll.Poll.ID,
		&Poll{
			ChatID:         chatID,
			WordLineNumber: correctLineNumber,
			Language:       "sh",
			CreateEpoch:    time.Now().Unix(),
		},
	)
	if err != nil {
		log.Printf("Error while saving poll: %s", err)
	}
}

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

	// Using the Config value, create the DynamoDB client
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

	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Debug = true

	s3BucketId := os.Getenv("s3_bucket_id")
	log.Printf("s3BucketId = %s\n", s3BucketId)
	s3Client := s3.NewFromConfig(cfg)
	dictionary, err := loadDictionary()
	if err != nil {
		log.Printf("Error while loading dictionary: %s", err)
		return
	}

	if update.Message != nil { // If we got a message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		sendPoll(ctx, s3Client, s3BucketId, dictionary, bot, update.Message.Chat.ID)

	} else if update.Poll != nil {
		poll, err := GetPoll(ctx, s3Client, s3BucketId, update.Poll.ID)
		if err != nil {
			log.Printf("Error while retrieving poll: %s", err)
		}

		sendPoll(ctx, s3Client, s3BucketId, dictionary, bot, poll.ChatID)
	}
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Body size = %d. \n", len(request.Body))
	log.Printf("Body = %s. \n", request.Body)
	log.Println("Headers:")
	for key, value := range request.Headers {
		log.Printf("  %s: %s\n", key, value)
	}

	Reply(ctx, request.Body)
	return events.APIGatewayProxyResponse{Body: "POST", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
