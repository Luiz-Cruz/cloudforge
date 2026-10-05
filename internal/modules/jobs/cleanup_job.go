package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/services/notification"
	"github.com/Luiz-Cruz/cloudforge/internal/services/storage"
	"github.com/sirupsen/logrus"
)

type cleanupFailedJob struct {
	envRepo      environment.Repository
	storage      storage.Service
	notification notification.Service
	maxAge       time.Duration
}

func NewCleanupFailedJob(
	repo environment.Repository,
	storageService storage.Service,
	notificationService notification.Service,
	maxAge time.Duration,
) Job {
	if maxAge <= 0 {
		maxAge = 24 * time.Hour
	}
	return &cleanupFailedJob{
		envRepo:      repo,
		storage:      storageService,
		notification: notificationService,
		maxAge:       maxAge,
	}
}

func (j *cleanupFailedJob) RunJob(ctx context.Context) error {
	logrus.Infof("Executing CleanupFailedJob: scanning for environments FAILED older than %v", j.maxAge)

	envs, err := j.envRepo.FindAll(ctx)
	if err != nil {
		logrus.Errorf("CleanupFailedJob failed to list environments: %v", err)
		return err
	}

	threshold := time.Now().UTC().Add(-j.maxAge)
	archivedCount := 0

	for _, env := range envs {
		if env.Status == "FAILED" && env.CreatedAt.Before(threshold) {
			logrus.Infof("Archiving stale failed environment: %s (Created: %v)", env.TransactionID, env.CreatedAt)

			if err := j.envRepo.UpdateStatus(ctx, env.TransactionID, "ARCHIVED"); err != nil {
				logrus.Errorf("Failed to update status to ARCHIVED for %s: %v", env.TransactionID, err)
				continue
			}

			archivedCount++
		}
	}

	logrus.Infof("CleanupFailedJob completed: %d environments archived", archivedCount)

	if archivedCount > 0 && j.notification != nil {
		subject := "CloudForge Cron: Maintenance Cleanup Completed"
		msg := fmt.Sprintf("Scheduled maintenance archived %d stale failed environments.", archivedCount)
		if err := j.notification.NotifyTopic(ctx, subject, msg); err != nil {
			logrus.Errorf("Failed to publish cleanup notification: %v", err)
		}
	}

	return nil
}
