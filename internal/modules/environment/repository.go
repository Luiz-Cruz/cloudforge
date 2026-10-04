package environment

import (
	"context"

	"github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (*EnvironmentState, error)
	FindAll(ctx context.Context) ([]*EnvironmentState, error)
	SaveState(ctx context.Context, state EnvironmentState) error
	UpdateStatus(ctx context.Context, transactionID, status string) error
}

type dynamoRepository struct {
	client    wrapper.DynamoDBAPI
	tableName string
}

func NewRepository(client wrapper.DynamoDBAPI, tableName string) Repository {
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

func (r *dynamoRepository) FindByID(ctx context.Context, id string) (*EnvironmentState, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"transaction_id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil || len(out.Item) == 0 {
		return nil, nil
	}

	var state EnvironmentState
	err = attributevalue.UnmarshalMap(out.Item, &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *dynamoRepository) FindAll(ctx context.Context) ([]*EnvironmentState, error) {
	out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &r.tableName,
	})
	if err != nil {
		return nil, err
	}

	var states []*EnvironmentState
	err = attributevalue.UnmarshalListOfMaps(out.Items, &states)
	if err != nil {
		return nil, err
	}
	return states, nil
}
