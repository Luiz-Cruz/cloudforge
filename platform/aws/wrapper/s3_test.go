package wrapper

import (
	"context"
	"errors"
	"testing"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	mock_wrapper "github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestS3Wrapper_PutObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const fileName = "manifest.json"
	const bucket = "cloudforge-artifacts"
	data := []byte(`{"status":"AVAILABLE"}`)

	api := mock_wrapper.NewMockS3API(ctrl)
	wrapper := NewS3APIWrapper(api, bucket, nil)
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		api.EXPECT().PutObject(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(ctx context.Context,
			params *s3.PutObjectInput,
			optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			assert.Equal(t, bucket, *params.Bucket)
			assert.Equal(t, fileName, *params.Key)
			assert.NotNil(t, params.Body)
			return &s3.PutObjectOutput{}, nil
		})
		assert.NoError(t, wrapper.PutObject(ctx, fileName, data))
	})

	t.Run("Fail", func(t *testing.T) {
		api.EXPECT().PutObject(gomock.Eq(ctx), gomock.Any()).Return(nil, errors.New("s3 upload failed"))
		assert.Error(t, wrapper.PutObject(ctx, fileName, data))
	})
}

func TestS3Wrapper_PresignGetObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const fileName = "manifest.json"
	const bucket = "cloudforge-artifacts"
	const signedURL = "https://cloudforge.s3.amazonaws.com/manifest.json?sig=xyz"

	api := mock_wrapper.NewMockS3SignedAPI(ctrl)
	wrapper := NewS3APIWrapper(nil, bucket, api)
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		api.EXPECT().PresignGetObject(gomock.Eq(ctx), gomock.Any()).Return(&v4.PresignedHTTPRequest{
			URL: signedURL,
		}, nil)

		result, err := wrapper.PresignGetObject(ctx, fileName)
		assert.NoError(t, err)
		assert.Equal(t, signedURL, result)
	})

	t.Run("Nil SignedAPI", func(t *testing.T) {
		wrapperNil := NewS3APIWrapper(nil, bucket, nil)
		res, err := wrapperNil.PresignGetObject(ctx, fileName)
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("Fail", func(t *testing.T) {
		api.EXPECT().PresignGetObject(gomock.Eq(ctx), gomock.Any()).Return(nil, errors.New("presign failed"))
		_, err := wrapper.PresignGetObject(ctx, fileName)
		assert.Error(t, err)
	})
}

func TestS3Wrapper_HeadObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const fileName = "manifest.json"
	const bucket = "cloudforge-artifacts"

	api := mock_wrapper.NewMockS3API(ctrl)
	wrapper := NewS3APIWrapper(api, bucket, nil)
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		cType := "application/json"
		size := int64(128)

		api.EXPECT().HeadObject(gomock.Eq(ctx), gomock.Any()).Return(&s3.HeadObjectOutput{
			LastModified:  &now,
			ContentType:   &cType,
			ContentLength: &size,
		}, nil)

		result, err := wrapper.HeadObject(ctx, fileName)
		assert.NoError(t, err)
		assert.Equal(t, &now, result["LastModified"])
		assert.Equal(t, cType, result["ContentType"])
		assert.Equal(t, size, result["Size"])
	})

	t.Run("Fail", func(t *testing.T) {
		api.EXPECT().HeadObject(gomock.Eq(ctx), gomock.Any()).Return(nil, errors.New("head failed"))
		_, err := wrapper.HeadObject(ctx, fileName)
		assert.Error(t, err)
	})
}
