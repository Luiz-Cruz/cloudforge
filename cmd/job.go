package cmd

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

type JobApplication struct{}

type JobEvent struct {
	Name string `json:"name"`
}

func (JobApplication) Run() {
	lambda.Start(func(ctx context.Context, event JobEvent) error {
		logrus.Infof("Received EventBridge Job execution: %v", event)
		// TODO: Implement job routing based on event.Name
		return nil
	})
}
