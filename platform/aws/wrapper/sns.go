package wrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

//go:generate mockgen -source=./sns.go -destination=./mock/sns_mock.go -package=mock
type SNSAPI interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}
