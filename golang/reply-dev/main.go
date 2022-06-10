package main

import (
	"context"

	telegram "example.com/telegram"
)

func main() {
	ctx := context.TODO()
	requestBody := "{\"poll\":{    \"id\": \"5846093446460211733\",    \"question\": \"peacemaker, pacifier (Noun)\",    \"options\": [        {            \"text\": \"vijećnik\",            \"voter_count\": 0        },        {            \"text\": \"mirotvorac\",            \"voter_count\": 1        },        {            \"text\": \"riznica\",            \"voter_count\": 0        }    ],    \"total_voter_count\": 1,    \"is_closed\": false,    \"is_anonymous\": true,    \"type\": \"quiz\",    \"allows_multiple_answers\": false,    \"correct_option_id\": 1}}"
	telegram.Reply(ctx, requestBody)
}
