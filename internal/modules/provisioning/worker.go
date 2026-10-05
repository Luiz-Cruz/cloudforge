package provisioning

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
)

type Worker interface {
	ProcessSQS(ctx context.Context, sqsEvent events.SQSEvent) error
}

type ProvisioningStep struct {
	Name     string
	Execute  func(ctx context.Context, state environment.EnvironmentState) error
	Rollback func(ctx context.Context, state environment.EnvironmentState) error
}

type provisioningWorker struct {
	envRepo environment.Repository
	steps   []ProvisioningStep
	tracer  wrapper.Tracer
}

func DefaultSteps() []ProvisioningStep {
	return []ProvisioningStep{
		{
			Name: "Network",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Infof("[Tx: %s] Provisioning Network...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
			Rollback: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Warnf("[Tx: %s] Rolling back Network...", state.TransactionID)
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
				logrus.Warnf("[Tx: %s] Rolling back Database...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
		},
		{
			Name: "Storage",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Infof("[Tx: %s] Provisioning Storage...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
			Rollback: func(ctx context.Context, state environment.EnvironmentState) error {
				logrus.Warnf("[Tx: %s] Rolling back Storage...", state.TransactionID)
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
				logrus.Warnf("[Tx: %s] Rolling back Compute...", state.TransactionID)
				time.Sleep(200 * time.Millisecond)
				return nil
			},
		},
	}
}

func NewWorker(repo environment.Repository, tracer wrapper.Tracer) Worker {
	return NewWorkerWithSteps(repo, DefaultSteps(), tracer)
}

func NewWorkerWithSteps(repo environment.Repository, steps []ProvisioningStep, tracer wrapper.Tracer) Worker {
	if tracer == nil {
		tracer = wrapper.NewNoopTracer()
	}
	return &provisioningWorker{
		envRepo: repo,
		steps:   steps,
		tracer:  tracer,
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

		if err := w.executeProvisioning(ctx, state); err != nil {
			logrus.Errorf("Provisioning failed for %s. Rollback completed.", state.TransactionID)
			w.envRepo.UpdateStatus(ctx, state.TransactionID, "FAILED")
			return err
		}

		logrus.Infof("Successfully provisioned environment %s", state.TransactionID)
		w.envRepo.UpdateStatus(ctx, state.TransactionID, "AVAILABLE")
	}

	return nil
}

func (w *provisioningWorker) executeProvisioning(ctx context.Context, state environment.EnvironmentState) error {
	var successfulSteps []ProvisioningStep

	for _, step := range w.steps {
		stepName := step.Name
		err := w.tracer.Capture(ctx, fmt.Sprintf("ProvisioningStep:%s", stepName), func(stepCtx context.Context) error {
			return step.Execute(stepCtx, state)
		})

		if err != nil {
			logrus.Errorf("[Tx: %s] Error executing %s: %v", state.TransactionID, stepName, err)

			// Execute Compensating Transactions (Rollback) in Reverse Order with subsegment tracing
			for i := len(successfulSteps) - 1; i >= 0; i-- {
				rbStep := successfulSteps[i]
				rbName := rbStep.Name
				_ = w.tracer.Capture(ctx, fmt.Sprintf("RollbackStep:%s", rbName), func(rbCtx context.Context) error {
					rbErr := rbStep.Rollback(rbCtx, state)
					if rbErr != nil {
						logrus.Errorf("[Tx: %s] FATAL: Rollback failed for %s: %v", state.TransactionID, rbName, rbErr)
					}
					return rbErr
				})
			}
			return err
		}
		successfulSteps = append(successfulSteps, step)
	}

	return nil
}
