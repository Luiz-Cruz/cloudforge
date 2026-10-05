package environment

import "time"

type EnvironmentState struct {
	TransactionID string    `json:"transaction_id" dynamodbav:"transaction_id"`
	Name          string    `json:"name" dynamodbav:"name"`
	Type          string    `json:"type" dynamodbav:"type"`
	Status        string    `json:"status" dynamodbav:"status"`
	CreatedAt     time.Time `json:"created_at" dynamodbav:"created_at"`
}
