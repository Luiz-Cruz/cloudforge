package cdi

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	appConfig "github.com/Luiz-Cruz/cloudforge/platform/config"
)

func LoadAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	return cfg
}

func ProvideDynamoDB() *dynamodb.Client {
	cfg := LoadAWSConfig()
	if appConfig.IsLocalStack() {
		return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String("http://127.0.0.1:4566")
		})
	}
	return dynamodb.NewFromConfig(cfg)
}

func ProvideSQS() *sqs.Client {
	cfg := LoadAWSConfig()
	if appConfig.IsLocalStack() {
		return sqs.NewFromConfig(cfg, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String("http://127.0.0.1:4566")
		})
	}
	return sqs.NewFromConfig(cfg)
}

func ProvideSNS() *sns.Client {
	cfg := LoadAWSConfig()
	if appConfig.IsLocalStack() {
		return sns.NewFromConfig(cfg, func(o *sns.Options) {
			o.BaseEndpoint = aws.String("http://127.0.0.1:4566")
		})
	}
	return sns.NewFromConfig(cfg)
}

func ProvideS3() *s3.Client {
	cfg := LoadAWSConfig()
	if appConfig.IsLocalStack() {
		return s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String("http://127.0.0.1:4566")
		})
	}
	return s3.NewFromConfig(cfg)
}
