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
	if len(os.Args) < 2 {
		log.Println("Usage: ./main <telegramBotURLs>")
		os.Exit(1)
	}

	var telegramBotURLs map[string]string
	err := json.Unmarshal([]byte(os.Args[1]), &telegramBotURLs)
	if err != nil {
		log.Printf("unable to parse telegram bot tokens, %v", err)
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return
	}

	telegramBotTokens, err := telegram.GetTokens(ctx, cfg)
	if err != nil {
		log.Printf("unable to get telegram bot tokens, %v", err)
	}

	for language, token := range telegramBotTokens {
		telegramBotURL := telegramBotURLs[language]
		log.Printf("language = %s, url = %s\n", language, telegramBotURL)
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Authorized on account %s", bot.Self.UserName)

		// Register device with tgbotapi
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(telegramBotURL))
		if err != nil {
			log.Panic(err)
		}

		webHookInfo, err := bot.GetWebhookInfo()
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Webhook info: %+v\n", webHookInfo)
	}
}
