package provisioning_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment/mocks"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/provisioning"
	"github.com/Luiz-Cruz/cloudforge/internal/services/email"
	mock_email "github.com/Luiz-Cruz/cloudforge/internal/services/email/mock"
	mock_notification "github.com/Luiz-Cruz/cloudforge/internal/services/notification/mock"
	mock_storage "github.com/Luiz-Cruz/cloudforge/internal/services/storage/mock"
	mock_wrapper "github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper/mock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWorker_ProcessSQS_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockTracer := mock_wrapper.NewMockTracer(ctrl)
	mockStorage := mock_storage.NewMockService(ctrl)
	mockNotification := mock_notification.NewMockService(ctrl)
	mockEmail := mock_email.NewMockService(ctrl)

	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-100", "PROCESSING").Return(nil).Times(1)
	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-100", "AVAILABLE").Return(nil).Times(1)

	// Tracer captures every step
	mockTracer.EXPECT().Capture(gomock.Any(), "ProvisioningStep:StepA", gomock.Any()).
		DoAndReturn(func(ctx context.Context, name string, fn func(context.Context) error) error {
			return fn(ctx)
		}).Times(1)

	mockTracer.EXPECT().Capture(gomock.Any(), "ProvisioningStep:StepB", gomock.Any()).
		DoAndReturn(func(ctx context.Context, name string, fn func(context.Context) error) error {
			return fn(ctx)
		}).Times(1)

	// Storage expectation
	mockStorage.EXPECT().SaveFile(gomock.Any(), "manifests/tx-100.json", gomock.Any()).Return(nil).Times(1)
	mockStorage.EXPECT().GetFileURL(gomock.Any(), "manifests/tx-100.json").Return("https://s3.local/manifests/tx-100.json", nil).Times(1)

	// Notification expectation
	mockNotification.EXPECT().NotifyTopic(gomock.Any(), "Environment Available: prod-app", gomock.Any()).Return(nil).Times(1)

	// Email expectation
	mockEmail.EXPECT().SendTemplateEmail(
		gomock.Any(),
		email.EnvironmentProvisionedTemplate,
		gomock.Any(),
		"CloudForge: Environment prod-app Ready",
		"admin@cloudforge.io",
	).Return(nil).Times(1)

	executedSteps := []string{}
	customSteps := []provisioning.ProvisioningStep{
		{
			Name: "StepA",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				executedSteps = append(executedSteps, "StepA")
				return nil
			},
		},
		{
			Name: "StepB",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				executedSteps = append(executedSteps, "StepB")
				return nil
			},
		},
	}

	worker := provisioning.NewWorkerWithSteps(mockRepo, customSteps, mockTracer, mockStorage, mockEmail, mockNotification)

	state := environment.EnvironmentState{
		TransactionID: "tx-100",
		Name:          "prod-app",
		Type:          "kubernetes",
	}
	body, _ := json.Marshal(state)

	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{MessageId: "msg-1", Body: string(body)},
		},
	}

	err := worker.ProcessSQS(context.Background(), event)
	assert.NoError(t, err)
	assert.Equal(t, []string{"StepA", "StepB"}, executedSteps)
}

func TestWorker_ProcessSQS_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	worker := provisioning.NewWorker(mockRepo, nil, nil, nil, nil)

	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{MessageId: "msg-invalid", Body: "{broken-json"},
		},
	}

	err := worker.ProcessSQS(context.Background(), event)
	assert.NoError(t, err)
}

func TestWorker_ProcessSQS_StepFailureAndRollback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockTracer := mock_wrapper.NewMockTracer(ctrl)
	mockNotification := mock_notification.NewMockService(ctrl)
	mockEmail := mock_email.NewMockService(ctrl)

	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-200", "PROCESSING").Return(nil).Times(1)
	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-200", "FAILED").Return(nil).Times(1)

	expectedErr := errors.New("provisioning step B failed")

	mockTracer.EXPECT().Capture(gomock.Any(), "ProvisioningStep:Step1", gomock.Any()).
		DoAndReturn(func(ctx context.Context, name string, fn func(context.Context) error) error {
			return fn(ctx)
		}).Times(1)

	mockTracer.EXPECT().Capture(gomock.Any(), "ProvisioningStep:Step2", gomock.Any()).
		DoAndReturn(func(ctx context.Context, name string, fn func(context.Context) error) error {
			return fn(ctx)
		}).Times(1)

	mockTracer.EXPECT().Capture(gomock.Any(), "RollbackStep:Step1", gomock.Any()).
		DoAndReturn(func(ctx context.Context, name string, fn func(context.Context) error) error {
			return fn(ctx)
		}).Times(1)

	mockNotification.EXPECT().NotifyTopic(gomock.Any(), "Environment Failed: fail-app", gomock.Any()).Return(nil).Times(1)

	mockEmail.EXPECT().SendTemplateEmail(
		gomock.Any(),
		email.EnvironmentFailedTemplate,
		gomock.Any(),
		"CloudForge Alert: Provisioning Failed for fail-app",
		"admin@cloudforge.io",
	).Return(nil).Times(1)

	rollbackCalled := false

	customSteps := []provisioning.ProvisioningStep{
		{
			Name: "Step1",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				return nil
			},
			Rollback: func(ctx context.Context, state environment.EnvironmentState) error {
				rollbackCalled = true
				return nil
			},
		},
		{
			Name: "Step2",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				return expectedErr
			},
			Rollback: func(ctx context.Context, state environment.EnvironmentState) error {
				return nil
			},
		},
	}

	worker := provisioning.NewWorkerWithSteps(mockRepo, customSteps, mockTracer, nil, mockEmail, mockNotification)

	state := environment.EnvironmentState{
		TransactionID: "tx-200",
		Name:          "fail-app",
		Type:          "database",
	}
	body, _ := json.Marshal(state)

	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{MessageId: "msg-fail", Body: string(body)},
		},
	}

	err := worker.ProcessSQS(context.Background(), event)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.True(t, rollbackCalled)
}

func TestWorker_ProcessSQS_NilServices_Graceful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-300", "PROCESSING").Return(nil).Times(1)
	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-300", "AVAILABLE").Return(nil).Times(1)

	customSteps := []provisioning.ProvisioningStep{
		{
			Name: "StepA",
			Execute: func(ctx context.Context, state environment.EnvironmentState) error {
				return nil
			},
		},
	}

	worker := provisioning.NewWorkerWithSteps(mockRepo, customSteps, nil, nil, nil, nil)

	state := environment.EnvironmentState{
		TransactionID: "tx-300",
		Name:          "app-nil-services",
	}
	body, _ := json.Marshal(state)

	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{MessageId: "msg-300", Body: string(body)},
		},
	}

	err := worker.ProcessSQS(context.Background(), event)
	assert.NoError(t, err)
}

func TestWorker_DefaultSteps(t *testing.T) {
	steps := provisioning.DefaultSteps()
	assert.Len(t, steps, 4)

	state := environment.EnvironmentState{TransactionID: "tx-default"}
	ctx := context.Background()

	for _, s := range steps {
		assert.NotEmpty(t, s.Name)
		err := s.Execute(ctx, state)
		assert.NoError(t, err)
		err = s.Rollback(ctx, state)
		assert.NoError(t, err)
	}
}
