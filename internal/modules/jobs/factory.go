package jobs

import (
	"time"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/services/notification"
	"github.com/Luiz-Cruz/cloudforge/internal/services/storage"
)

type Factory struct {
	envRepo      environment.Repository
	storage      storage.Service
	notification notification.Service
}

func NewFactory(repo environment.Repository, storageService storage.Service, notificationService notification.Service) *Factory {
	return &Factory{
		envRepo:      repo,
		storage:      storageService,
		notification: notificationService,
	}
}

func (f *Factory) BuildJob(jobType JobType) Job {
	switch jobType {
	case CleanupFailedJob, ArchiveStaleJob:
		return NewCleanupFailedJob(f.envRepo, f.storage, f.notification, 24*time.Hour)
	default:
		return nil
	}
}
