package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Bucket interface {
	DownloadImage(ctx context.Context, sess *session.Session, bucket string, key string) (bool, []byte, error)
	UploadImage(ctx context.Context, sess *session.Session, bucket string, key string, data []byte) (*s3manager.UploadOutput, error)
}

type S3Handler struct{}

func (s *S3Handler) DownloadImage(ctx context.Context, sess *session.Session, bucket string, key string) (bool, []byte, error) {
	//Check existence of object
	//Validate s3 object
	srv := s3.New(sess)
	_, err := srv.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return false, nil, nil
	}
	//Download image from s3
	buffer := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.DownloadWithContext(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return true, nil, errors.New(fmt.Sprintf("could not download image from S3: %v", err))
	}
	return true, buffer.Bytes(), nil
}

func (s *S3Handler) UploadImage(ctx context.Context, sess *session.Session, bucket string, key string, data []byte) (*s3manager.UploadOutput, error) {
	uploader := s3manager.NewUploader(sess)
	output, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String("image/jpeg"),
		Body:        bytes.NewReader(data),
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}
