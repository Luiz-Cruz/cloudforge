package cmd

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/provisioning"
	"github.com/Luiz-Cruz/cloudforge/platform/cdi"
)

type JobApplication struct{}

func (JobApplication) Run() {
	dynamoClient := cdi.ProvideDynamoDB()
	tableName := "cloudforge-saga-state-local"
	
	repo := environment.NewRepository(dynamoClient, tableName)
	worker := provisioning.NewWorker(repo)

	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) error {
		logrus.Infof("Job woke up with %d SQS records", len(sqsEvent.Records))
		return worker.ProcessSQS(ctx, sqsEvent)
	})
}
