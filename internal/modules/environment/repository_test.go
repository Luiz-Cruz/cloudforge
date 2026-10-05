package environment_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	mock_wrapper "github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper/mock"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRepository_SaveState_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := environment.NewRepository(mockDynamo, "test-table")

	mockDynamo.EXPECT().PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil).Times(1)

	state := environment.EnvironmentState{
		TransactionID: "tx-save-1",
		Name:          "save-app",
		Status:        "PENDING",
		CreatedAt:     time.Now().UTC(),
	}

	err := repo.SaveState(context.Background(), state)
	assert.NoError(t, err)
}

func TestRepository_SaveState_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := environment.NewRepository(mockDynamo, "test-table")

	expectedErr := errors.New("dynamo put failed")
	mockDynamo.EXPECT().PutItem(gomock.Any(), gomock.Any()).
		Return(nil, expectedErr).Times(1)

	state := environment.EnvironmentState{TransactionID: "tx-save-err"}
	err := repo.SaveState(context.Background(), state)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestRepository_UpdateStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := environment.NewRepository(mockDynamo, "test-table")

	mockDynamo.EXPECT().UpdateItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.UpdateItemOutput{}, nil).Times(1)

	err := repo.UpdateStatus(context.Background(), "tx-up-1", "AVAILABLE")
	assert.NoError(t, err)
}

func TestRepository_FindByID_Found(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := environment.NewRepository(mockDynamo, "test-table")

	expectedState := environment.EnvironmentState{
		TransactionID: "tx-find-1",
		Name:          "found-app",
		Status:        "AVAILABLE",
	}
	item, _ := attributevalue.MarshalMap(expectedState)

	mockDynamo.EXPECT().GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil).Times(1)

	state, err := repo.FindByID(context.Background(), "tx-find-1")
	assert.NoError(t, err)
	assert.NotNil(t, state)
	assert.Equal(t, "tx-find-1", state.TransactionID)
}

func TestRepository_FindByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := environment.NewRepository(mockDynamo, "test-table")

	mockDynamo.EXPECT().GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: nil}, nil).Times(1)

	state, err := repo.FindByID(context.Background(), "tx-notfound")
	assert.NoError(t, err)
	assert.Nil(t, state)
}

func TestRepository_FindByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := environment.NewRepository(mockDynamo, "test-table")

	mockDynamo.EXPECT().GetItem(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("dynamo read error")).Times(1)

	state, err := repo.FindByID(context.Background(), "tx-err")
	assert.Error(t, err)
	assert.Nil(t, state)
}

func TestRepository_FindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := environment.NewRepository(mockDynamo, "test-table")

	item1, _ := attributevalue.MarshalMap(environment.EnvironmentState{TransactionID: "tx-1", Name: "app1"})
	item2, _ := attributevalue.MarshalMap(environment.EnvironmentState{TransactionID: "tx-2", Name: "app2"})

	mockDynamo.EXPECT().Scan(gomock.Any(), gomock.Any()).
		Return(&dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{item1, item2},
		}, nil).Times(1)

	states, err := repo.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, states, 2)
}

func TestRepository_FindAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := mock_wrapper.NewMockDynamoDBAPI(ctrl)
	repo := environment.NewRepository(mockDynamo, "test-table")

	mockDynamo.EXPECT().Scan(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("scan error")).Times(1)

	states, err := repo.FindAll(context.Background())
	assert.Error(t, err)
	assert.Nil(t, states)
}
