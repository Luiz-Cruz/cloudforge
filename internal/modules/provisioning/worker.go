package provisioning

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
)

type Worker interface {
	ProcessSQS(ctx context.Context, sqsEvent events.SQSEvent) error
}

type provisioningWorker struct {
	envRepo environment.Repository
}

func NewWorker(repo environment.Repository) Worker {
	return &provisioningWorker{
		envRepo: repo,
	}
}

func (w *provisioningWorker) ProcessSQS(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		logrus.Infof("Processing SQS message ID: %s", message.MessageId)

		var state environment.EnvironmentState
		if err := json.Unmarshal([]byte(message.Body), &state); err != nil {
			logrus.Errorf("Failed to unmarshal SQS message body: %v", err)
			continue
		}

		if err := w.simulateProvisioning(ctx, state); err != nil {
			logrus.Errorf("Provisioning failed for %s. Executing Rollback.", state.TransactionID)
			w.envRepo.UpdateStatus(ctx, state.TransactionID, "FAILED")
			return err
		}

		logrus.Infof("Successfully provisioned environment %s", state.TransactionID)
		w.envRepo.UpdateStatus(ctx, state.TransactionID, "AVAILABLE")
	}

	return nil
}

func (w *provisioningWorker) simulateProvisioning(ctx context.Context, state environment.EnvironmentState) error {
	steps := []string{"Network", "Database", "Storage", "Compute", "Secrets"}

	for _, step := range steps {
		logrus.Infof("[Tx: %s] Provisioning %s...", state.TransactionID, step)
		time.Sleep(500 * time.Millisecond) // Simulate work

	}

	return nil
}
