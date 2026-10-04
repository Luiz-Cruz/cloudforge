package provisioning

import (
	"context"
	"encoding/json"
	"errors"
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

		// Update state to PROCESSING
		w.envRepo.UpdateStatus(ctx, state.TransactionID, "PROCESSING")

		if err := w.executeSaga(ctx, state); err != nil {
			logrus.Errorf("Provisioning failed for %s. Rollback completed.", state.TransactionID)
			w.envRepo.UpdateStatus(ctx, state.TransactionID, "FAILED")
			// Return error so the message goes back to queue or DLQ
			return err
		}

		logrus.Infof("Successfully provisioned environment %s", state.TransactionID)
		w.envRepo.UpdateStatus(ctx, state.TransactionID, "AVAILABLE")
	}

	return nil
}

type SagaStep struct {
	Name     string
	Execute  func(ctx context.Context, state environment.EnvironmentState) error
	Rollback func(ctx context.Context, state environment.EnvironmentState) error
}

func (w *provisioningWorker) executeSaga(ctx context.Context, state environment.EnvironmentState) error {
	steps := []SagaStep{
		{
			Name: "Network",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Infof("[Tx: %s] Provisioning Network...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
			Rollback: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Warnf("[Tx: %s] ⏪ Rolling back Network...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
		},
		{
			Name: "Database",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Infof("[Tx: %s] Provisioning Database...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
			Rollback: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Warnf("[Tx: %s] ⏪ Rolling back Database...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
		},
		{
			Name: "Storage",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Infof("[Tx: %s] Provisioning Storage...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				// SIMULATING A RANDOM FAILURE HERE IF THE NAME CONTAINS "fail"
				if state.Name == "fail" {
					return errors.New("simulated storage failure")
				}
				return nil
			},
			Rollback: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Warnf("[Tx: %s] ⏪ Rolling back Storage...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
		},
		{
			Name: "Compute",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Infof("[Tx: %s] Provisioning Compute...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
			Rollback: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Warnf("[Tx: %s] ⏪ Rolling back Compute...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
		},
	}

	var successfulSteps []SagaStep

	for _, step := range steps {
		err := step.Execute(ctx, state)
		if err != nil {
			logrus.Errorf("[Tx: %s] ❌ Error executing %s: %v", state.TransactionID, step.Name, err)

			// Execute Compensating Transactions (Rollback) in Reverse Order
			for i := len(successfulSteps) - 1; i >= 0; i-- {
				rbStep := successfulSteps[i]
				rbErr := rbStep.Rollback(ctx, state)
				if rbErr != nil {
					logrus.Errorf("[Tx: %s] ❌ FATAL: Rollback failed for %s: %v", state.TransactionID, rbStep.Name, rbErr)
				}
			}
			return err
		}
		successfulSteps = append(successfulSteps, step)
	}

	return nil
}
