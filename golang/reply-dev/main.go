package main

import (
	"context"

	telegram "example.com/telegram"
)

func main() {
	ctx := context.Background()
	requestBody := `{
		"update_id": 9629202,
		"poll": {
		  "id": "5846093446460212749",
		  "question": "peacemaker, pacifier (Noun)",
		  "options": [
			{
			  "text": "vijećnik",
			  "voter_count": 0
			},
			{
			  "text": "mirotvorac",
			  "voter_count": 1
			},
			{
			  "text": "riznica",
			  "voter_count": 0
			}
		  ],
		  "total_voter_count": 1,
		  "is_closed": false,
		  "is_anonymous": true,
		  "type": "quiz",
		  "allows_multiple_answers": false,
		  "correct_option_id": 1
		}
	  }`
	s3BucketId := "my-bucket-legible-quetzal"
	languageCode := "sh"
	telegram.Reply(ctx, requestBody, s3BucketId, languageCode)
	languageCode = "nl"
	telegram.Reply(ctx, requestBody, s3BucketId, languageCode)
}
