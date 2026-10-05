package cdi

import (
	"context"

	"github.com/Luiz-Cruz/cloudforge/platform/aws/wrapper"
	appConfig "github.com/Luiz-Cruz/cloudforge/platform/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func LoadAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	return cfg
}

func provideDynamoDB() *dynamodb.Client {
	cfg := LoadAWSConfig()
	if appConfig.IsLocalStack() {
		return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String("http://127.0.0.1:4566")
		})
	}
	return dynamodb.NewFromConfig(cfg)
}

func provideSQS() wrapper.SQSAPI {
	cfg := LoadAWSConfig()
	if appConfig.IsLocalStack() {
		return sqs.NewFromConfig(cfg, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String("http://127.0.0.1:4566")
		})
	}
	return sqs.NewFromConfig(cfg)
}

func provideSNS() *sns.Client {
	cfg := LoadAWSConfig()
	if appConfig.IsLocalStack() {
		return sns.NewFromConfig(cfg, func(o *sns.Options) {
			o.BaseEndpoint = aws.String("http://127.0.0.1:4566")
		})
	}
	return sns.NewFromConfig(cfg)
}

func provideS3() *s3.Client {
	cfg := LoadAWSConfig()
	if appConfig.IsLocalStack() {
		return s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String("http://127.0.0.1:4566")
		})
	}
	return s3.NewFromConfig(cfg)
}

func provideSES() wrapper.SESAPI {
	if appConfig.IsLocalStack() {
		return wrapper.NewMockSESAPI()
	}
	cfg := LoadAWSConfig()
	return sesv2.NewFromConfig(cfg)
}

func ProvideDynamoDB() *dynamodb.Client { return provideDynamoDB() }
func ProvideSQS() wrapper.SQSAPI        { return provideSQS() }
func ProvideSNS() *sns.Client           { return provideSNS() }
func ProvideS3() *s3.Client             { return provideS3() }
func ProvideSES() wrapper.SESAPI        { return provideSES() }

func provideTracer() wrapper.Tracer {
	if appConfig.IsLocalStack() {
		return wrapper.NewNoopTracer()
	}
	return wrapper.NewXRayTracer()
}

func ProvideTracer() wrapper.Tracer { return provideTracer() }
