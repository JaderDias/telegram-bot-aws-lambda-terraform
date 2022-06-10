package telegram

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type MockS3Client struct {
	PutObjectInput *s3.PutObjectInput
	GetObjectInput *s3.GetObjectInput
}

func (c *MockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	c.PutObjectInput = params
	return nil, nil
}

func (c *MockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	c.GetObjectInput = params
	if params.Key == nil {
		return nil, fmt.Errorf("key is nil")
	}

	result := "{\"Languages\":{\"sh\":{\"RightAnswers\":{\"1\":true,\"2\":true,\"3\":true},\"WrongAnswers\":{\"4\":true,\"5\":true,\"6\":true}}}}"
	if *params.Key == "poll/1" {
		result = "{\"ChatID\":123,\"WordId\":456,\"Language\":\"sh\"}"
	} else if *params.Key == "poll/0" {
		result = ""
	}

	return &s3.GetObjectOutput{
		Body:          io.NopCloser(strings.NewReader(result)),
		ContentLength: int64(len(result)),
	}, nil
}

func PutObject(
	ctx context.Context,
	s3Client S3Client,
	s3BucketId string,
	key string,
	value []byte,
	expires *time.Time,
) error {
	log.Printf("PutObject: %s", key)
	_, err := s3Client.PutObject(ctx,
		&s3.PutObjectInput{
			Body:    bytes.NewReader(value),
			Bucket:  &s3BucketId,
			Key:     &key,
			Expires: expires,
		})
	return err
}

func ToString(reader io.Reader) string {
	var bbuffer bytes.Buffer
	_, err := io.Copy(&bbuffer, reader)
	if err != nil {
		return ""
	}
	return bbuffer.String()
}

func GetObject(
	ctx context.Context,
	s3Client S3Client,
	s3BucketId string,
	key string,
) ([]byte, error) {
	log.Printf("GetObject: %s", key)
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s3BucketId,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.ContentLength == 0 {
		return nil, fmt.Errorf("contentLength is 0")
	}

	log.Printf("ContentLength %d", resp.ContentLength)
	var bbuffer bytes.Buffer
	buffer := make([]byte, resp.ContentLength)
	for {
		num, rerr := resp.Body.Read(buffer)
		if num > 0 {
			bbuffer.Write(buffer[:num])
		} else if rerr == io.EOF || rerr != nil {
			break
		}
	}
	return bbuffer.Bytes(), nil
}
