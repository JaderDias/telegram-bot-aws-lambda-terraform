package main

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func PutObject(
	ctx context.Context,
	s3Client S3Client,
	s3BucketId string,
	key string,
	value []byte,
) error {
	expires := time.Now().Add(time.Hour * 24 * 30) // 30 days
	_, err := s3Client.PutObject(ctx,
		&s3.PutObjectInput{
			Body:    bytes.NewReader(value),
			Bucket:  &s3BucketId,
			Key:     &key,
			Expires: &expires,
		})
	return err
}

func GetObject(
	ctx context.Context,
	s3Client *s3.Client,
	s3BucketId string,
	key string,
) (bytes.Buffer, error) {
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s3BucketId,
		Key:    &key,
	})
	defer resp.Body.Close()
	var bbuffer bytes.Buffer
	if err != nil {
		return bbuffer, err
	}

	buffer := make([]byte, resp.ContentLength)
	for {
		num, rerr := resp.Body.Read(buffer)
		if num > 0 {
			bbuffer.Write(buffer[:num])
		} else if rerr == io.EOF || rerr != nil {
			break
		}
	}
	return bbuffer, nil
}
