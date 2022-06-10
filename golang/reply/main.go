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

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	HTTPMethodNotSupported = errors.New("HTTP method not supported")
)

func loadDictionary() []string {
	file, err := os.Open("sh.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return words
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
	dictionary []string,
	bot *tgbotapi.BotAPI,
	chatID int64,
) {
	correctLineNumber := -1
	correctLineNumber, sendPollConfig := getPoll(dictionary, correctLineNumber)
	sendPollConfig.BaseChat = tgbotapi.BaseChat{
		ChatID: chatID,
	}

	_, err := bot.Send(sendPollConfig)
	if err != nil {
		log.Printf("Error while sending poll: %s", err)
		return
	}
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Body size = %d. \n", len(request.Body))
	fmt.Printf("Body = %s. \n", request.Body)
	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}

	var update tgbotapi.Update
	err := json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		log.Println(err)
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	ssmClient := ssm.NewFromConfig(cfg)

	parameterName := "telegram_bot_token"
	parameterOutput, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &parameterName,
		WithDecryption: true,
	})
	if err != nil {
		log.Fatalf("unable to get telegram bot token, %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(*parameterOutput.Parameter.Value)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	dictionary := loadDictionary()

	if update.Message != nil { // If we got a message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		sendPoll(dictionary, bot, update.Message.Chat.ID)

	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return events.APIGatewayProxyResponse{Body: "POST", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
