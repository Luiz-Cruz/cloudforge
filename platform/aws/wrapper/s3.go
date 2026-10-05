package wrapper

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./s3.go -destination=./mock/s3_mock.go -package=mock
type S3API interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

type S3SignedAPI interface {
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type S3APIWrapper interface {
	PutObject(ctx context.Context, key string, data []byte) error
	PresignGetObject(ctx context.Context, key string) (string, error)
	HeadObject(ctx context.Context, key string) (map[string]interface{}, error)
}

type s3Wrapper struct {
	api       S3API
	signedAPI S3SignedAPI
	bucket    string
}

func NewS3APIWrapper(api S3API, bucket string, signedAPI S3SignedAPI) S3APIWrapper {
	return &s3Wrapper{
		api:       api,
		signedAPI: signedAPI,
		bucket:    bucket,
	}
}

func (w *s3Wrapper) PutObject(ctx context.Context, key string, data []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(w.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}

	_, err := w.api.PutObject(ctx, input)
	if err != nil {
		logrus.Errorf("Error uploading object to S3 bucket: %v", err)
		return err
	}
	return nil
}

func (w *s3Wrapper) PresignGetObject(ctx context.Context, key string) (string, error) {
	if w.signedAPI == nil {
		return "", nil
	}
	input := &s3.GetObjectInput{
		Bucket: aws.String(w.bucket),
		Key:    aws.String(key),
	}

	request, err := w.signedAPI.PresignGetObject(ctx, input)
	if err != nil {
		logrus.Errorf("Error generating signed url: %v", err)
		return "", err
	}
	return request.URL, nil
}

func (w *s3Wrapper) HeadObject(ctx context.Context, key string) (map[string]interface{}, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(w.bucket),
		Key:    aws.String(key),
	}

	response, err := w.api.HeadObject(ctx, input)
	if err != nil {
		logrus.Errorf("Error heading object from S3 bucket: %v", err)
		return nil, err
	}

	metadata := make(map[string]interface{})
	metadata["LastModified"] = response.LastModified
	if response.ContentType != nil {
		metadata["ContentType"] = *response.ContentType
	}
	if response.ContentLength != nil {
		metadata["Size"] = *response.ContentLength
	}

	return metadata, nil
}
