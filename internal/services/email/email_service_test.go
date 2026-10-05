package email

import (
	"context"
	"errors"
	"testing"

	mock_wrapper "github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper/mock"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	testSubject   = "CloudForge Notification"
	testReceiver1 = "devops@example.com"
	testReceiver2 = "lead@example.com"
)

func TestEmailService_SendTemplateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sesMock := mock_wrapper.NewMockSESAPI(ctrl)
	service := NewEmailService(sesMock, "no-reply@cloudforge.io")

	data := map[string]string{
		"Title": "Test Environment",
		"Body":  "Environment ready for deployment",
	}

	t.Run("Given a valid template, when sent successfully, then returns no error", func(t *testing.T) {
		sesMock.EXPECT().
			SendEmail(gomock.Any(), gomock.Any()).
			Return(&sesv2.SendEmailOutput{}, nil)

		err := service.SendTemplateEmail(context.Background(), TestTemplate, data, testSubject, testReceiver1, testReceiver2)
		assert.NoError(t, err)
	})

	t.Run("Given an invalid template name, then returns error", func(t *testing.T) {
		err := service.SendTemplateEmail(context.Background(), "non_existent_template", data, testSubject, testReceiver1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "template non_existent_template not found")
	})

	t.Run("Given a valid template, when SES returns error, then returns error", func(t *testing.T) {
		sesMock.EXPECT().
			SendEmail(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("ses quota exceeded"))

		err := service.SendTemplateEmail(context.Background(), TestTemplate, data, testSubject, testReceiver1)
		assert.Error(t, err)
		assert.Equal(t, "ses quota exceeded", err.Error())
	})

	t.Run("Given an environment provisioned template, when sent, renders and calls SES", func(t *testing.T) {
		provData := map[string]string{
			"TransactionID":   "tx-12345",
			"EnvironmentName": "production-eu",
			"ManifestURL":     "https://s3.amazonaws.com/cloudforge/manifest.json",
		}
		sesMock.EXPECT().
			SendEmail(gomock.Any(), gomock.Any()).
			Return(&sesv2.SendEmailOutput{}, nil)

		err := service.SendTemplateEmail(context.Background(), EnvironmentProvisionedTemplate, provData, "Provisioned", testReceiver1)
		assert.NoError(t, err)
	})

	t.Run("Given an environment failed template, when sent, renders and calls SES", func(t *testing.T) {
		failData := map[string]string{
			"TransactionID":   "tx-99999",
			"EnvironmentName": "staging-us",
			"Error":           "network allocation timeout",
		}
		sesMock.EXPECT().
			SendEmail(gomock.Any(), gomock.Any()).
			Return(&sesv2.SendEmailOutput{}, nil)

		err := service.SendTemplateEmail(context.Background(), EnvironmentFailedTemplate, failData, "Failed", testReceiver1)
		assert.NoError(t, err)
	})
}

func TestEmailService_BuildEmailInput(t *testing.T) {
	const emailBody = "<html><body>Hello CloudForge</body></html>"
	const sender = "noreply@cloudforge.io"

	emailInput := buildEmailInput(sender, testSubject, emailBody, testReceiver1, testReceiver2)

	assert.Equal(t, sender, *emailInput.FromEmailAddress)
	assert.Len(t, emailInput.Destination.ToAddresses, 2)
	assert.Equal(t, testReceiver1, emailInput.Destination.ToAddresses[0])
	assert.Equal(t, testReceiver2, emailInput.Destination.ToAddresses[1])
	assert.NotNil(t, emailInput.Content.Simple.Subject)
	assert.Equal(t, testSubject, *emailInput.Content.Simple.Subject.Data)
	assert.NotNil(t, emailInput.Content.Simple.Body)
	assert.Equal(t, emailBody, *emailInput.Content.Simple.Body.Html.Data)
}
