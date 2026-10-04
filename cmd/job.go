package cmd

import (
	"context"

	"github.com/Luiz-Cruz/cloudforge/platform/cdi"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type JobApplication struct{}

func (JobApplication) Run() {
	worker := cdi.ProvideProvisioningWorker()

	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) error {
		logrus.Infof("Job woke up with %d SQS records", len(sqsEvent.Records))
		return worker.ProcessSQS(ctx, sqsEvent)
	})
}

type LocalJobApplication struct{}

func (LocalJobApplication) Run() {
	worker := cdi.ProvideProvisioningWorker()
	sqsClient := cdi.ProvideSQS()

	queueURL := viper.GetString("QUEUE_URL")
	if queueURL == "" {
		queueURL = "http://localhost:4566/000000000000/cloudforge-queue-local"
	}

	logrus.Info("Starting local SQS Poller for Job Worker...")
	for {
		msgResult, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
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
				for _, m := range msgResult.Messages {
					sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
						QueueUrl:      aws.String(queueURL),
						ReceiptHandle: m.ReceiptHandle,
					})
				}
				logrus.Info("Batch processed successfully")
			}
		}
	}
}
