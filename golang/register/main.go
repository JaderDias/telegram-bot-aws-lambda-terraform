package main

import (
	"log"
	"os"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	telegramBotToken := os.Args[1]
	telegramBotURL := os.Args[2]
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Register device with tgbotapi
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(telegramBotURL))
	if err != nil {
		log.Panic(err)
	}
}
