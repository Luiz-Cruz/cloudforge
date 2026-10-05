package notification

import (
	"context"

	"github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./notification.go -destination=./mock/notification_mock.go -package=mock
type Service interface {
	NotifyTopic(ctx context.Context, subject, text string) error
}

type snsService struct {
	api   wrapper.SNSAPI
	topic string
}

func NewService(api wrapper.SNSAPI, topic string) Service {
	return &snsService{
		api:   api,
		topic: topic,
	}
}

func (s *snsService) NotifyTopic(ctx context.Context, subject, text string) error {
	logrus.Infof("Publishing notification to SNS topic %s: %s", s.topic, subject)
	_, err := s.api.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(s.topic),
		Subject:  aws.String(subject),
		Message:  aws.String(text),
	})
	return err
}
