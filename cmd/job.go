package cmd

import (
	"context"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/provisioning"
	"github.com/Luiz-Cruz/cloudforge/platform/cdi"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
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

type LocalJobApplication struct{}

func (LocalJobApplication) Run() {
	dynamoClient := cdi.ProvideDynamoDB()
	tableName := "cloudforge-saga-state-local"
	repo := environment.NewRepository(dynamoClient, tableName)
	worker := provisioning.NewWorker(repo)

	sqsClient := cdi.ProvideSQS()
	queueUrl := "http://localhost:4566/000000000000/cloudforge-queue-local"

	logrus.Info("Starting local SQS Poller for Job Worker...")
	for {
		msgResult, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueUrl),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     5,
		})
		if err != nil {
			logrus.Errorf("Error receiving messages: %v", err)
			continue
		}

		if len(msgResult.Messages) > 0 {
			var sqsRecords []events.SQSMessage
			for _, m := range msgResult.Messages {
				sqsRecords = append(sqsRecords, events.SQSMessage{
					MessageId: *m.MessageId,
					Body:      *m.Body,
				})
			}

			event := events.SQSEvent{Records: sqsRecords}
			logrus.Infof("Polled %d messages, processing...", len(sqsRecords))
			err = worker.ProcessSQS(context.TODO(), event)
			if err != nil {
				logrus.Errorf("Error processing batch: %v", err)
			} else {
				// Delete processed messages
				for _, m := range msgResult.Messages {
					sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
						QueueUrl:      aws.String(queueUrl),
						ReceiptHandle: m.ReceiptHandle,
					})
				}
				logrus.Info("Batch processed successfully")
			}
		}
	}
}
