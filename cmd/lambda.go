package cmd

import (
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
)

type LambdaApplication struct{}

func (LambdaApplication) Run() {
	app := provideFiberApplication()
	fiberLambda := fiberadapter.New(app)
	lambda.Start(fiberLambda.ProxyWithContext)
}
