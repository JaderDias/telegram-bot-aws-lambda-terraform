package main

import (
	"context"
	"log"
	"os"

	telegram "example.com/telegram"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("Usage: ./main <tokenParameterName> <s3BucketId>")
		os.Exit(1)
	}

	tokenParameterName := os.Args[1]
	s3BucketId := os.Args[2]
	ctx := context.Background()
	requestBody := `{"update_id":9629336,"message":{"message_id":7025,"from":{"id":5299480268,"is_bot":false,"first_name":"J","last_name":"D","language_code":"en"},"chat":{"id":5299480268,"first_name":"J","last_name":"D","type":"private"},"date":1655282459,"text":"/word","entities":[{"offset":0,"length":5,"type":"bot_command"}]}}`
	languageCode := "nl"
	telegram.Reply(ctx, requestBody, s3BucketId, languageCode, tokenParameterName)
}
