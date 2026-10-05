package notification_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Luiz-Cruz/cloudforge/internal/services/notification"
	mock_wrapper "github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper/mock"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotification_NotifyTopic_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSNS := mock_wrapper.NewMockSNSAPI(ctrl)
	svc := notification.NewService(mockSNS, "arn:aws:sns:us-east-1:123456789:cloudforge-topic")

	mockSNS.EXPECT().Publish(gomock.Any(), gomock.Any()).
		Return(&sns.PublishOutput{}, nil).Times(1)

	err := svc.NotifyTopic(context.Background(), "Environment Ready", "Environment is available")
	assert.NoError(t, err)
}

func TestNotification_NotifyTopic_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSNS := mock_wrapper.NewMockSNSAPI(ctrl)
	svc := notification.NewService(mockSNS, "arn:aws:sns:us-east-1:123456789:cloudforge-topic")

	mockSNS.EXPECT().Publish(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("sns error")).Times(1)

	err := svc.NotifyTopic(context.Background(), "Environment Ready", "Environment is available")
	assert.Error(t, err)
}
