package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetBot(tokens map[string]string, language string) (*tgbotapi.BotAPI, error) {
	log.Printf("GetBot language = %s\n", language)
	bot, err := tgbotapi.NewBotAPI(tokens[language])
	if err != nil {
		log.Printf("Error while creating bot: %s", err)
		return nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Debug = true
	return bot, nil
}
