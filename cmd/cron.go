package cmd

import (
	"context"
	"errors"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/jobs"
	"github.com/Luiz-Cruz/cloudforge/platform/cdi"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CronEvent struct {
	Name string `json:"name"`
}

type CronApplication struct{}

func (CronApplication) Run() {
	factory := cdi.ProvideJobsFactory()

	lambda.Start(func(ctx context.Context, event CronEvent) error {
		logrus.Infof("Cron execution started with event: %+v", event)
		jobName := event.Name
		if jobName == "" {
			jobName = string(jobs.CleanupFailedJob)
		}

		var jobType jobs.JobType
		jt, err := jobType.From(jobName)
		if err != nil {
			logrus.Errorf("Invalid cron job name: %s", jobName)
			return err
		}

		job := factory.BuildJob(jt)
		if job == nil {
			return errors.New("job not found")
		}

		return job.RunJob(ctx)
	})
}

type LocalCronApplication struct{}

func (LocalCronApplication) Run() {
	logrus.Info("Starting local Cron job execution...")
	factory := cdi.ProvideJobsFactory()

	jobName := viper.GetString("JOB_NAME")
	if jobName == "" {
		jobName = string(jobs.CleanupFailedJob)
	}

	var jobType jobs.JobType
	jt, err := jobType.From(jobName)
	if err != nil {
		logrus.Fatalf("Invalid job type: %s", jobName)
	}

	job := factory.BuildJob(jt)
	if job == nil {
		logrus.Fatalf("Job implementation not found for %s", jt)
	}

	if err := job.RunJob(context.Background()); err != nil {
		logrus.Fatalf("Cron job failed: %v", err)
	}
	logrus.Info("Local Cron execution finished successfully.")
}
