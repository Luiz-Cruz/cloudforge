package email

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"

	"github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

//go:generate mockgen -source=./email_service.go -destination=./mock/email_service_mock.go -package=mock
type Service interface {
	SendTemplateEmail(ctx context.Context, templateName string, data any, subject string, recipients ...string) error
}

//go:embed templates/*
var fs embed.FS

const (
	EnvironmentProvisionedTemplate = "environment_provisioned"
	EnvironmentFailedTemplate      = "environment_failed"
	TestTemplate                   = "test_template"
)

type emailService struct {
	ses       wrapper.SESAPI
	fromEmail string
	templates map[string]*template.Template
}

func NewEmailService(ses wrapper.SESAPI, fromEmail string) Service {
	templates := make(map[string]*template.Template)
	templates[EnvironmentProvisionedTemplate] = createTemplate(EnvironmentProvisionedTemplate, string(loadEmailHTML("templates/environment_provisioned.html")))
	templates[EnvironmentFailedTemplate] = createTemplate(EnvironmentFailedTemplate, string(loadEmailHTML("templates/environment_failed.html")))
	templates[TestTemplate] = createTemplate(TestTemplate, string(loadEmailHTML("templates/test.html")))
	return &emailService{
		ses:       ses,
		fromEmail: fromEmail,
		templates: templates,
	}
}

func (s *emailService) SendTemplateEmail(ctx context.Context, templateName string, data any, subject string, recipients ...string) error {
	tmpl, found := s.templates[templateName]
	if !found {
		return fmt.Errorf("template %s not found", templateName)
	}

	body, err := parseTemplate(tmpl, data)
	if err != nil {
		return err
	}

	input := buildEmailInput(s.fromEmail, subject, body, recipients...)
	_, err = s.ses.SendEmail(ctx, input)
	return err
}

func buildEmailInput(fromEmail, subject, body string, recipients ...string) *sesv2.SendEmailInput {
	return &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(body),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
		Destination: &types.Destination{
			ToAddresses: recipients,
		},
		FromEmailAddress: aws.String(fromEmail),
	}
}

func loadEmailHTML(fileName string) []byte {
	data, err := fs.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return data
}

func createTemplate(name, htmlContent string) *template.Template {
	return template.Must(template.New(name).Parse(htmlContent))
}

func parseTemplate(tmpl *template.Template, data any) (string, error) {
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
