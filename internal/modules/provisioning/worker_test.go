package provisioning_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment/mocks"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/provisioning"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWorker_ProcessSQS_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-100", "PROCESSING").Return(nil).Times(1)
	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-100", "AVAILABLE").Return(nil).Times(1)

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

	worker := provisioning.NewWorkerWithSteps(mockRepo, customSteps)

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
	worker := provisioning.NewWorker(mockRepo)

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

	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-200", "PROCESSING").Return(nil).Times(1)
	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-200", "FAILED").Return(nil).Times(1)

	rollbackCalled := false
	expectedErr := errors.New("provisioning step B failed")

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

	worker := provisioning.NewWorkerWithSteps(mockRepo, customSteps)

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

func TestWorker_DefaultSteps(t *testing.T) {
	steps := provisioning.DefaultSteps()
	assert.Len(t, steps, 4)

	state := environment.EnvironmentState{TransactionID: "tx-default"}
	ctx := context.Background()

	// Verify each default step runs and rolls back without error
	for _, s := range steps {
		assert.NotEmpty(t, s.Name)
		err := s.Execute(ctx, state)
		assert.NoError(t, err)
		err = s.Rollback(ctx, state)
		assert.NoError(t, err)
	}
}
