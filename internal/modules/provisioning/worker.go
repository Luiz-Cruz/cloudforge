package provisioning

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/services/email"
	"github.com/Luiz-Cruz/cloudforge/internal/services/notification"
	"github.com/Luiz-Cruz/cloudforge/internal/services/storage"
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

type EnvironmentManifest struct {
	TransactionID string    `json:"transaction_id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	CompletedAt   time.Time `json:"completed_at"`
}

type provisioningWorker struct {
	envRepo      environment.Repository
	steps        []ProvisioningStep
	tracer       wrapper.Tracer
	storage      storage.Service
	email        email.Service
	notification notification.Service
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

func NewWorker(
	repo environment.Repository,
	tracer wrapper.Tracer,
	storageService storage.Service,
	emailService email.Service,
	notificationService notification.Service,
) Worker {
	return NewWorkerWithSteps(repo, DefaultSteps(), tracer, storageService, emailService, notificationService)
}

func NewWorkerWithSteps(
	repo environment.Repository,
	steps []ProvisioningStep,
	tracer wrapper.Tracer,
	storageService storage.Service,
	emailService email.Service,
	notificationService notification.Service,
) Worker {
	if tracer == nil {
		tracer = wrapper.NewNoopTracer()
	}
	return &provisioningWorker{
		envRepo:      repo,
		steps:        steps,
		tracer:       tracer,
		storage:      storageService,
		email:        emailService,
		notification: notificationService,
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
			w.handleFailure(ctx, state, err)
			return err
		}

		logrus.Infof("Successfully provisioned environment %s", state.TransactionID)
		w.envRepo.UpdateStatus(ctx, state.TransactionID, "AVAILABLE")
		w.handleSuccess(ctx, state)
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

func (w *provisioningWorker) handleSuccess(ctx context.Context, state environment.EnvironmentState) {
	manifestURL := ""
	if w.storage != nil {
		manifest := EnvironmentManifest{
			TransactionID: state.TransactionID,
			Name:          state.Name,
			Type:          state.Type,
			Status:        "AVAILABLE",
			CompletedAt:   time.Now().UTC(),
		}
		manifestData, err := json.MarshalIndent(manifest, "", "  ")
		if err == nil {
			manifestKey := fmt.Sprintf("manifests/%s.json", state.TransactionID)
			if saveErr := w.storage.SaveFile(ctx, manifestKey, manifestData); saveErr != nil {
				logrus.Errorf("[Tx: %s] Failed to save manifest to storage: %v", state.TransactionID, saveErr)
			} else {
				if url, urlErr := w.storage.GetFileURL(ctx, manifestKey); urlErr == nil {
					manifestURL = url
				}
			}
		}
	}

	if w.notification != nil {
		subject := fmt.Sprintf("Environment Available: %s", state.Name)
		msg := fmt.Sprintf("Environment %s (Transaction: %s) has been successfully provisioned and is AVAILABLE.", state.Name, state.TransactionID)
		if err := w.notification.NotifyTopic(ctx, subject, msg); err != nil {
			logrus.Errorf("[Tx: %s] Failed to send SNS notification: %v", state.TransactionID, err)
		}
	}

	if w.email != nil {
		emailData := map[string]string{
			"TransactionID":   state.TransactionID,
			"EnvironmentName": state.Name,
			"ManifestURL":     manifestURL,
		}
		subject := fmt.Sprintf("CloudForge: Environment %s Ready", state.Name)
		if err := w.email.SendTemplateEmail(ctx, email.EnvironmentProvisionedTemplate, emailData, subject, "admin@cloudforge.io"); err != nil {
			logrus.Errorf("[Tx: %s] Failed to send SES email: %v", state.TransactionID, err)
		}
	}
}

func (w *provisioningWorker) handleFailure(ctx context.Context, state environment.EnvironmentState, provErr error) {
	if w.notification != nil {
		subject := fmt.Sprintf("Environment Failed: %s", state.Name)
		msg := fmt.Sprintf("Environment %s (Transaction: %s) failed provisioning: %v. Rollback completed.", state.Name, state.TransactionID, provErr)
		if err := w.notification.NotifyTopic(ctx, subject, msg); err != nil {
			logrus.Errorf("[Tx: %s] Failed to send SNS notification: %v", state.TransactionID, err)
		}
	}

	if w.email != nil {
		emailData := map[string]string{
			"TransactionID":   state.TransactionID,
			"EnvironmentName": state.Name,
			"Error":           provErr.Error(),
		}
		subject := fmt.Sprintf("CloudForge Alert: Provisioning Failed for %s", state.Name)
		if err := w.email.SendTemplateEmail(ctx, email.EnvironmentFailedTemplate, emailData, subject, "admin@cloudforge.io"); err != nil {
			logrus.Errorf("[Tx: %s] Failed to send SES email: %v", state.TransactionID, err)
		}
	}
}
