package provisioning_test

import (
	"context"
	"encoding/json"
	"testing"
	"github.com/aws/aws-lambda-go/events"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment/mocks"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/provisioning"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWorker_ProcessSQS_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	worker := provisioning.NewWorker(mockRepo)

	// Setup valid SQS message
	state := environment.EnvironmentState{
		TransactionID: "tx-123",
		Name:          "prod-db",
		Type:          "production",
		Status:        "PENDING",
	}
	bodyBytes, _ := json.Marshal(state)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: "msg-123",
				Body:      string(bodyBytes),
			},
		},
	}

	// Expect the repo to update status to AVAILABLE on success
	mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-123", "AVAILABLE").Return(nil).Times(1)

	err := worker.ProcessSQS(context.Background(), sqsEvent)
	
	assert.NoError(t, err)
}

func TestWorker_ProcessSQS_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	worker := provisioning.NewWorker(mockRepo)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: "msg-bad",
				Body:      "{invalid-json}",
			},
		},
	}

	// Repo shouldn't be called if JSON is invalid
	mockRepo.EXPECT().UpdateStatus(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	// It ignores invalid JSON to not block the queue and returns nil
	err := worker.ProcessSQS(context.Background(), sqsEvent)
	assert.NoError(t, err)
}
