package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	batchLines = 100
)

func uploadFile(
	ctx context.Context,
	s3Client *s3.Client,
	s3BucketId string,
	fileName string,
) {
	language := strings.TrimSuffix(fileName, ".csv")
	log.Printf("uploading %s", language)

	inFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	batchId := 0
	lines := make([]string, batchLines)
	lineIndex := 0

	for scanner.Scan() {
		lines[lineIndex] = scanner.Text()
		lineIndex++
		if lineIndex == batchLines {
			key := fmt.Sprintf("language/%s/%d", language, batchId)
			log.Printf("creating %s", key)
			_, err := s3Client.PutObject(
				ctx,
				&s3.PutObjectInput{
					Bucket: &s3BucketId,
					Key:    &key,
					Body:   strings.NewReader(strings.Join(lines, "\n")),
				},
			)
			if err != nil {
				log.Fatal(err)
			}
			batchId++
			lineIndex = 0
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func main() {
	if len(os.Args) < 3 {
		log.Panic("missing argument: <awsRegion> <s3BucketId>")
	}

	awsRegion := os.Args[1]
	s3BucketId := os.Args[2]

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return
	}

	s3Client := s3.NewFromConfig(cfg)
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".csv") {
			uploadFile(ctx, s3Client, s3BucketId, file.Name())
		}
	}
}
