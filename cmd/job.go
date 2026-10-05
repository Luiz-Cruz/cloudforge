package cmd

import (
	"context"
	"strconv"
	"time"

	"github.com/Luiz-Cruz/cloudforge/platform/cdi"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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

	dlqURL := viper.GetString("DLQ_URL")
	if dlqURL == "" {
		dlqURL = "http://localhost:4566/000000000000/cloudforge-dlq-local"
	}

	maxReceiveCount := viper.GetInt("MAX_RECEIVE_COUNT")
	if maxReceiveCount <= 0 {
		maxReceiveCount = 3
	}

	logrus.Info("Starting local SQS Poller with DLQ and Exponential Backoff resilience...")
	for {
		msgResult, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:                    aws.String(queueURL),
			MaxNumberOfMessages:         10,
			WaitTimeSeconds:             5,
			MessageSystemAttributeNames: []types.MessageSystemAttributeName{types.MessageSystemAttributeNameApproximateReceiveCount},
		})
		if err != nil {
			logrus.Errorf("Error receiving messages: %v", err)
			time.Sleep(2 * time.Second)
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
				logrus.Errorf("Batch processing failed: %v", err)
				for _, m := range msgResult.Messages {
					receiveCount := 1
					if countStr, ok := m.Attributes[string(types.MessageSystemAttributeNameApproximateReceiveCount)]; ok {
						if parsed, parseErr := strconv.Atoi(countStr); parseErr == nil {
							receiveCount = parsed
						}
					}

					if receiveCount >= maxReceiveCount {
						logrus.Warnf("Message %s exceeded max retries (%d/%d). Moving to DLQ: %s", *m.MessageId, receiveCount, maxReceiveCount, dlqURL)
						_, dlqErr := sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
							QueueUrl:    aws.String(dlqURL),
							MessageBody: m.Body,
						})
						if dlqErr != nil {
							logrus.Errorf("Failed to forward message %s to DLQ: %v", *m.MessageId, dlqErr)
						} else {
							sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
								QueueUrl:      aws.String(queueURL),
								ReceiptHandle: m.ReceiptHandle,
							})
							logrus.Infof("Message %s routed to DLQ and deleted from primary queue", *m.MessageId)
						}
					} else {
						backoffDelay := time.Duration(1<<receiveCount) * 100 * time.Millisecond
						logrus.Warnf("Retrying message %s (attempt %d). Backing off for %v", *m.MessageId, receiveCount, backoffDelay)
						time.Sleep(backoffDelay)
					}
				}
			} else {
				for _, m := range msgResult.Messages {
					sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
						QueueUrl:      aws.String(queueURL),
						ReceiptHandle: m.ReceiptHandle,
					})
				}
				logrus.Info("Batch processed and evicted from queue successfully")
			}
		}
	}
}
