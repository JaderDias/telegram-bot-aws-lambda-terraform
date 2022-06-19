package telegram

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BaseChatParams(chat *tgbotapi.BaseChat) (tgbotapi.Params, error) {
	params := make(tgbotapi.Params)

	params.AddFirstValid("chat_id", chat.ChatID, chat.ChannelUsername)
	params.AddNonZero("reply_to_message_id", chat.ReplyToMessageID)
	params.AddBool("disable_notification", chat.DisableNotification)
	params.AddBool("allow_sending_without_reply", chat.AllowSendingWithoutReply)
	//	params.AddBool("protect_content", chat.ProtectContent)

	err := params.AddInterface("reply_markup", chat.ReplyMarkup)

	return params, err
}

func Params(config tgbotapi.SendPollConfig) (tgbotapi.Params, error) {
	params, err := BaseChatParams(&config.BaseChat)
	if err != nil {
		return params, err
	}

	params["question"] = config.Question
	if err = params.AddInterface("options", config.Options); err != nil {
		return params, err
	}
	params["is_anonymous"] = strconv.FormatBool(config.IsAnonymous)
	params.AddNonEmpty("type", config.Type)
	params["allows_multiple_answers"] = strconv.FormatBool(config.AllowsMultipleAnswers)
	params["correct_option_id"] = strconv.FormatInt(config.CorrectOptionID, 10)
	params.AddBool("is_closed", config.IsClosed)
	params.AddNonEmpty("explanation", config.Explanation)
	params.AddNonEmpty("explanation_parse_mode", config.ExplanationParseMode)
	params.AddNonZero("open_period", config.OpenPeriod)
	params.AddNonZero("close_date", config.CloseDate)
	err = params.AddInterface("explanation_entities", config.ExplanationEntities)

	return params, err
}
