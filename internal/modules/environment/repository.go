package environment

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

type Repository interface {
	SaveState(ctx context.Context, state EnvironmentState) error
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
