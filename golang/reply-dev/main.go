package main

import (
	"context"
	"log"
	"os"

	telegram "example.com/telegram"
)

type MockOsFileReader struct {
}

func (r *MockOsFileReader) ReadFile(fileName string) ([]byte, error) {
	return os.ReadFile("../../sh.csv")
}

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: ./main <tokenParameterName>")
		os.Exit(1)
	}

	tokenParameterName := os.Args[1]
	ctx := context.Background()
	requestBody := `{"update_id":9629336,"message":{"message_id":7025,"from":{"id":5299480268,"is_bot":false,"first_name":"J","last_name":"D","language_code":"en"},"chat":{"id":5299480268,"first_name":"J","last_name":"D","type":"private"},"date":1655282459,"text":"/word","entities":[{"offset":0,"length":5,"type":"bot_command"}]}}`
	languageCode := "sh"
	fileReader := &MockOsFileReader{}
	telegram.Reply(ctx, fileReader, requestBody, languageCode, tokenParameterName)
	languageCode = "nl"
	telegram.Reply(ctx, fileReader, requestBody, languageCode, tokenParameterName)
}
