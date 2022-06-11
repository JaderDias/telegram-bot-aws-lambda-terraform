package telegram

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetTokens(ctx context.Context, cfg aws.Config) (map[string]string, error) {
	ssmClient := ssm.NewFromConfig(cfg)

	parameterName := "telegram_bot_tokens"
	parameterOutput, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &parameterName,
		WithDecryption: true,
	})
	if err != nil {
		log.Printf("unable to get telegram bot tokens, %v", err)
		return nil, err
	}

	var telegramBotTokens map[string]string
	err = json.Unmarshal([]byte(*parameterOutput.Parameter.Value), &telegramBotTokens)
	return telegramBotTokens, err
}
