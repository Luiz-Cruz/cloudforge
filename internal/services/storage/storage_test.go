package storage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Luiz-Cruz/cloudforge/internal/services/storage"
	mock_wrapper "github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStorage_SaveFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mock_wrapper.NewMockS3APIWrapper(ctrl)
	svc := storage.NewS3Storage(mockS3)

	mockS3.EXPECT().PutObject(gomock.Any(), "manifests/tx-1.json", []byte("data")).
		Return(nil).Times(1)

	err := svc.SaveFile(context.Background(), "manifests/tx-1.json", []byte("data"))
	assert.NoError(t, err)
}

func TestStorage_SaveFile_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mock_wrapper.NewMockS3APIWrapper(ctrl)
	svc := storage.NewS3Storage(mockS3)

	mockS3.EXPECT().PutObject(gomock.Any(), "manifests/tx-1.json", []byte("data")).
		Return(errors.New("upload failed")).Times(1)

	err := svc.SaveFile(context.Background(), "manifests/tx-1.json", []byte("data"))
	assert.Error(t, err)
}

func TestStorage_GetFileURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mock_wrapper.NewMockS3APIWrapper(ctrl)
	svc := storage.NewS3Storage(mockS3)

	mockS3.EXPECT().PresignGetObject(gomock.Any(), "manifests/tx-1.json").
		Return("https://s3.amazonaws.com/presigned-url", nil).Times(1)

	url, err := svc.GetFileURL(context.Background(), "manifests/tx-1.json")
	assert.NoError(t, err)
	assert.Equal(t, "https://s3.amazonaws.com/presigned-url", url)
}
