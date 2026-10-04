package environment_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment/mocks"
	mock_wrapper "github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper/mock"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_StartProvisioning_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockSQS := mock_wrapper.NewMockSQSAPI(ctrl)

	mockRepo.EXPECT().SaveState(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockSQS.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(&sqs.SendMessageOutput{}, nil).Times(1)

	svc := environment.NewService(mockRepo, mockSQS, "http://dummy-queue")
	state, err := svc.StartProvisioning(context.Background(), "prod-db", "database")

	assert.NoError(t, err)
	assert.Equal(t, "prod-db", state.Name)
	assert.Equal(t, "database", state.Type)
	assert.Equal(t, "PENDING", state.Status)
	assert.NotEmpty(t, state.TransactionID)
}

func TestService_StartProvisioning_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockSQS := mock_wrapper.NewMockSQSAPI(ctrl)

	expectedErr := errors.New("dynamodb error")
	mockRepo.EXPECT().SaveState(gomock.Any(), gomock.Any()).Return(expectedErr).Times(1)

	svc := environment.NewService(mockRepo, mockSQS, "http://dummy-queue")
	state, err := svc.StartProvisioning(context.Background(), "prod-db", "database")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, state.TransactionID)
}

func TestService_StartProvisioning_SQSError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockSQS := mock_wrapper.NewMockSQSAPI(ctrl)

	mockRepo.EXPECT().SaveState(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	expectedErr := errors.New("sqs unreachable")
	mockSQS.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil, expectedErr).Times(1)

	svc := environment.NewService(mockRepo, mockSQS, "http://dummy-queue")
	state, err := svc.StartProvisioning(context.Background(), "prod-db", "database")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, state.TransactionID)
}

func TestService_GetEnvironment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockSQS := mock_wrapper.NewMockSQSAPI(ctrl)

	expectedState := &environment.EnvironmentState{
		TransactionID: "tx-123",
		Name:          "stage-app",
		Status:        "AVAILABLE",
	}
	mockRepo.EXPECT().FindByID(gomock.Any(), "tx-123").Return(expectedState, nil).Times(1)

	svc := environment.NewService(mockRepo, mockSQS, "http://dummy-queue")
	res, err := svc.GetEnvironment(context.Background(), "tx-123")

	assert.NoError(t, err)
	assert.Equal(t, expectedState, res)
}

func TestService_GetEnvironment_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockSQS := mock_wrapper.NewMockSQSAPI(ctrl)

	expectedErr := errors.New("not found")
	mockRepo.EXPECT().FindByID(gomock.Any(), "tx-missing").Return(nil, expectedErr).Times(1)

	svc := environment.NewService(mockRepo, mockSQS, "http://dummy-queue")
	res, err := svc.GetEnvironment(context.Background(), "tx-missing")

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestService_ListEnvironments_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockSQS := mock_wrapper.NewMockSQSAPI(ctrl)

	expectedList := []*environment.EnvironmentState{
		{TransactionID: "tx-1", Name: "env-1", Status: "AVAILABLE"},
		{TransactionID: "tx-2", Name: "env-2", Status: "PENDING"},
	}
	mockRepo.EXPECT().FindAll(gomock.Any()).Return(expectedList, nil).Times(1)

	svc := environment.NewService(mockRepo, mockSQS, "http://dummy-queue")
	list, err := svc.ListEnvironments(context.Background())

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, expectedList, list)
}

func TestService_ListEnvironments_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockSQS := mock_wrapper.NewMockSQSAPI(ctrl)

	expectedErr := errors.New("scan failed")
	mockRepo.EXPECT().FindAll(gomock.Any()).Return(nil, expectedErr).Times(1)

	svc := environment.NewService(mockRepo, mockSQS, "http://dummy-queue")
	list, err := svc.ListEnvironments(context.Background())

	assert.Error(t, err)
	assert.Nil(t, list)
}
