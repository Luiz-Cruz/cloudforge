package environment

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type Service interface {
	StartProvisioning(ctx context.Context, name, envType string) (EnvironmentState, error)
}

type environmentService struct {
	repo     Repository
	sqsCli   *sqs.Client
	queueUrl string
}

func NewService(repo Repository, sqsCli *sqs.Client, queueUrl string) Service {
	return &environmentService{
		repo:     repo,
		sqsCli:   sqsCli,
		queueUrl: queueUrl,
	}
}

func (s *environmentService) StartProvisioning(ctx context.Context, name, envType string) (EnvironmentState, error) {
	state := EnvironmentState{
		TransactionID: uuid.New().String(),
		Name:          name,
		Type:          envType,
		Status:        "PENDING",
		CreatedAt:     time.Now().UTC(),
	}

	if err := s.repo.SaveState(ctx, state); err != nil {
		return EnvironmentState{}, err
	}

	messageBody, _ := json.Marshal(state)
	_, err := s.sqsCli.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueUrl),
		MessageBody: aws.String(string(messageBody)),
	})

	if err != nil {
		logrus.Errorf("Failed to send SQS message: %v", err)
		return EnvironmentState{}, err
	}

	return state, nil
}
