package environment

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
	"time"
)

type EnvironmentState struct {
	TransactionID string    `dynamodbav:"transaction_id"`
	Name          string    `dynamodbav:"name"`
	Type          string    `dynamodbav:"type"`
	Status        string    `dynamodbav:"status"`
	CreatedAt     time.Time `dynamodbav:"created_at"`
}

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks
type Repository interface {
	SaveState(ctx context.Context, state EnvironmentState) error
	UpdateStatus(ctx context.Context, transactionID, status string) error
}

type dynamoRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) Repository {
	return &dynamoRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *dynamoRepository) SaveState(ctx context.Context, state EnvironmentState) error {
	item, err := attributevalue.MarshalMap(state)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		logrus.Errorf("Failed to save state in DynamoDB: %v", err)
		return err
	}

	return nil
}

func (r *dynamoRepository) UpdateStatus(ctx context.Context, transactionID, status string) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"transaction_id": &types.AttributeValueMemberS{Value: transactionID},
		},
		UpdateExpression: aws.String("SET #s = :status"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: status},
		},
	})
	if err != nil {
		logrus.Errorf("Failed to update status %s for tx %s: %v", status, transactionID, err)
	}
	return err
}
