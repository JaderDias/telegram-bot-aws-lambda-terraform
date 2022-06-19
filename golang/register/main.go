package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	telegram "example.com/telegram"
	"github.com/aws/aws-sdk-go-v2/config"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	if len(os.Args) < 4 {
		log.Println("Usage: ./main <awsRegion> <tokenParameterName> <telegramBotURLs>")
		os.Exit(1)
	}

	awsRegion := os.Args[1]
	tokenParameterName := os.Args[2]

	var telegramBotURLs map[string]string
	err := json.Unmarshal([]byte(os.Args[3]), &telegramBotURLs)
	if err != nil {
		log.Panicf("unable to parse telegram bot tokens, %v", err)
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		log.Panicf("unable to load SDK config, %v", err)
	}

	telegramBotTokens, err := telegram.GetTokens(ctx, cfg, tokenParameterName)
	if err != nil {
		log.Panicf("unable to get telegram bot tokens, %v", err)
	}

	for language, token := range telegramBotTokens {
		telegramBotURL := telegramBotURLs[language]
		log.Printf("language = %s, url = %s\n", language, telegramBotURL)
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Panicf("Error while creating bot: %s", err)
		}

		log.Printf("Authorized on account %s", bot.Self.UserName)

		// Register device with tgbotapi
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(telegramBotURL))
		if err != nil {
			log.Panicf("Error while setting webhook: %s", err)
		}

		webHookInfo, err := bot.GetWebhookInfo()
		if err != nil {
			log.Panicf("Error while getting webhook info: %s", err)
		}

		log.Printf("Webhook info: %+v\n", webHookInfo)
	}
}
