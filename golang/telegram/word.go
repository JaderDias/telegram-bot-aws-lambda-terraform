package telegram

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetWords(
	ctx context.Context,
	s3Client *s3.Client,
	s3BucketId string,
	language string,
	batchId int,
) ([]string, int, error) {
	if batchId == -1 {
		batchId = rand.Intn(260)
		log.Printf("random batchId: %d", batchId)
	}

	key := fmt.Sprintf("language/%s/%d", language, batchId)
	content, err := GetObject(ctx, s3Client, s3BucketId, key)
	if err != nil {
		return nil, 0, err
	}

	return strings.Split(string(content), "\n"), batchId, nil
}
