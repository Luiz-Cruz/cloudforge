package jobs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment"
	"github.com/Luiz-Cruz/cloudforge/internal/modules/environment/mocks"
	mock_notification "github.com/Luiz-Cruz/cloudforge/internal/services/notification/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCleanupFailedJob_RunJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockNotification := mock_notification.NewMockService(ctrl)

	now := time.Now().UTC()
	staleDate := now.Add(-48 * time.Hour)
	recentDate := now.Add(-2 * time.Hour)

	t.Run("Successfully archives stale failed environments and notifies", func(t *testing.T) {
		envs := []*environment.EnvironmentState{
			{
				TransactionID: "tx-stale-1",
				Status:        "FAILED",
				CreatedAt:     staleDate,
			},
			{
				TransactionID: "tx-recent-fail",
				Status:        "FAILED",
				CreatedAt:     recentDate,
			},
			{
				TransactionID: "tx-stale-available",
				Status:        "AVAILABLE",
				CreatedAt:     staleDate,
			},
		}

		mockRepo.EXPECT().FindAll(gomock.Any()).Return(envs, nil)
		mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-stale-1", "ARCHIVED").Return(nil)
		mockNotification.EXPECT().NotifyTopic(gomock.Any(), "CloudForge Cron: Maintenance Cleanup Completed", gomock.Any()).Return(nil)

		job := NewCleanupFailedJob(mockRepo, nil, mockNotification, 24*time.Hour)
		err := job.RunJob(context.Background())
		assert.NoError(t, err)
	})

	t.Run("Repository FindAll error returns error", func(t *testing.T) {
		mockRepo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("dynamo scan failed"))

		job := NewCleanupFailedJob(mockRepo, nil, nil, 0)
		err := job.RunJob(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "dynamo scan failed", err.Error())
	})

	t.Run("UpdateStatus error continues to next environment", func(t *testing.T) {
		envs := []*environment.EnvironmentState{
			{
				TransactionID: "tx-fail-update",
				Status:        "FAILED",
				CreatedAt:     staleDate,
			},
			{
				TransactionID: "tx-success-update",
				Status:        "FAILED",
				CreatedAt:     staleDate,
			},
		}

		mockRepo.EXPECT().FindAll(gomock.Any()).Return(envs, nil)
		mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-fail-update", "ARCHIVED").Return(errors.New("update err"))
		mockRepo.EXPECT().UpdateStatus(gomock.Any(), "tx-success-update", "ARCHIVED").Return(nil)
		mockNotification.EXPECT().NotifyTopic(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		job := NewCleanupFailedJob(mockRepo, nil, mockNotification, 24*time.Hour)
		err := job.RunJob(context.Background())
		assert.NoError(t, err)
	})
}
