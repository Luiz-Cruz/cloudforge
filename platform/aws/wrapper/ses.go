package wrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./ses.go -destination=./mock/ses_mock.go -package=mock
type SESAPI interface {
	SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

type mockSESAPI struct{}

func (m mockSESAPI) SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	logrus.Infof("Mocking email send: %v", params)
	return &sesv2.SendEmailOutput{}, nil
}

func NewMockSESAPI() SESAPI {
	return mockSESAPI{}
}
