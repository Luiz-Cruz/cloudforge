package environment

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

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
		TableName: &r.tableName,
		Item:      item,
	})

	return err
}

func (r *dynamoRepository) UpdateStatus(ctx context.Context, transactionID, status string) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.tableName,
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
	return err
}
