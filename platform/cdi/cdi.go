package cdi

import (
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/provisioning"
	"github.com/Luiz-Cruz/cloudforge/internal/services/email"
	"github.com/Luiz-Cruz/cloudforge/internal/services/notification"
	"github.com/Luiz-Cruz/cloudforge/internal/services/storage"
	"github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

var environmentService environment.Service
var environmentRepository environment.Repository
var provisioningWorker provisioning.Worker
var storageService storage.Service
var notificationService notification.Service
var emailService email.Service

func ProvideEnvironmentHandler() *environment.Handler {
	return environment.NewHandler(ProvideEnvironmentService())
}

func ProvideEnvironmentService() environment.Service {
	if environmentService == nil {
		queueURL := viper.GetString("QUEUE_URL")
		if queueURL == "" {
			queueURL = "http://localhost:4566/000000000000/cloudforge-queue-local"
		}
		environmentService = environment.NewService(
			provideEnvironmentRepository(),
			ProvideSQS(),
			queueURL,
		)
	}
	return environmentService
}

func provideEnvironmentRepository() environment.Repository {
	if environmentRepository == nil {
		tableName := viper.GetString("TABLE_NAME")
		if tableName == "" {
			tableName = "cloudforge-saga-state-local"
		}
		environmentRepository = environment.NewRepository(
			ProvideDynamoDB(),
			tableName,
		)
	}
	return environmentRepository
}

func ProvideStorageService() storage.Service {
	if storageService == nil {
		s3API := ProvideS3()
		s3SignedAPI := s3.NewPresignClient(s3API)
		bucket := viper.GetString("STORAGE_BUCKET")
		if bucket == "" {
			bucket = "cloudforge-artifacts-local"
		}
		storageService = storage.NewS3Storage(wrapper.NewS3APIWrapper(s3API, bucket, s3SignedAPI))
	}
	return storageService
}

func ProvideNotificationService() notification.Service {
	if notificationService == nil {
		topic := viper.GetString("NOTIFICATION_TOPIC_ARN")
		if topic == "" {
			topic = "arn:aws:sns:us-east-1:000000000000:cloudforge-topic-local"
		}
		notificationService = notification.NewService(ProvideSNS(), topic)
	}
	return notificationService
}

func ProvideEmailService() email.Service {
	if emailService == nil {
		sender := viper.GetString("EMAIL_SENDER")
		if sender == "" {
			sender = "no-reply@cloudforge.io"
		}
		emailService = email.NewEmailService(ProvideSES(), sender)
	}
	return emailService
}

func ProvideProvisioningWorker() provisioning.Worker {
	if provisioningWorker == nil {
		provisioningWorker = provisioning.NewWorker(
			provideEnvironmentRepository(),
			ProvideTracer(),
			ProvideStorageService(),
			ProvideEmailService(),
			ProvideNotificationService(),
		)
	}
	return provisioningWorker
}
