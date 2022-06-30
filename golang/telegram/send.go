package telegram

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Send(ctx context.Context, cfg aws.Config, s3BucketId, tokenParameterName string) {
	log.Printf("s3BucketId = %s\n", s3BucketId)
	s3Client := s3.NewFromConfig(cfg)
	prefix := "chat/"
	listObjectOut, err := s3Client.ListObjectsV2(
		ctx,
		&s3.ListObjectsV2Input{
			Bucket: &s3BucketId,
			Prefix: &prefix,
		},
	)
	if err != nil {
		log.Fatalf("unable to list objects, %v", err)
	}

	for _, obj := range listObjectOut.Contents {
		log.Printf("obj = %s\n", *obj.Key)
		chat, err := GetChatFromKey(
			ctx,
			s3Client,
			s3BucketId,
			*obj.Key,
		)
		if err != nil {
			log.Fatalf("unable to get chat, %v", err)
		}

		chatId, err := strconv.ParseInt(strings.Split(*obj.Key, "/")[1], 0, 64)
		if err != nil {
			log.Printf("unable to parse chatId, %v\n", err)
			continue
		}

		log.Printf("chat = %+v\n", chat)
		for languageCode, language := range chat.Languages {
			correctWordId := GetCorrectWordId(language)
			err = BotSendPollWithCorrectWordId(ctx, cfg, s3Client, s3BucketId, languageCode, chatId, tokenParameterName, correctWordId)
			if err != nil {
				log.Printf("Error while sending poll: %s", err)
			}
		}
	}
}
