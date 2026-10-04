package cdi

import (
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/provisioning"
	"github.com/spf13/viper"
)

var environmentService environment.Service
var environmentRepository environment.Repository
var provisioningWorker provisioning.Worker

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

func ProvideProvisioningWorker() provisioning.Worker {
	if provisioningWorker == nil {
		provisioningWorker = provisioning.NewWorker(provideEnvironmentRepository())
	}
	return provisioningWorker
}
