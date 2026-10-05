package storage

import (
	"context"

	"github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./storage.go -destination=./mock/storage_mock.go -package=mock
type Service interface {
	SaveFile(ctx context.Context, name string, data []byte) error
	GetFileURL(ctx context.Context, name string) (string, error)
}

type s3Storage struct {
	s3Wrapper wrapper.S3APIWrapper
}

func NewS3Storage(s3Wrapper wrapper.S3APIWrapper) Service {
	return &s3Storage{s3Wrapper: s3Wrapper}
}

func (s *s3Storage) SaveFile(ctx context.Context, name string, data []byte) error {
	logrus.Infof("Uploading file %s to storage", name)
	return s.s3Wrapper.PutObject(ctx, name, data)
}

func (s *s3Storage) GetFileURL(ctx context.Context, name string) (string, error) {
	return s.s3Wrapper.PresignGetObject(ctx, name)
}
