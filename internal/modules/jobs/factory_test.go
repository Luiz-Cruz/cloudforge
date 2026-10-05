package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory_BuildJob(t *testing.T) {
	factory := NewFactory(nil, nil, nil)

	t.Run("Build CleanupFailedJob", func(t *testing.T) {
		job := factory.BuildJob(CleanupFailedJob)
		assert.NotNil(t, job)
	})

	t.Run("Build ArchiveStaleJob", func(t *testing.T) {
		job := factory.BuildJob(ArchiveStaleJob)
		assert.NotNil(t, job)
	})

	t.Run("Build Unknown Job", func(t *testing.T) {
		job := factory.BuildJob(JobType("INVALID"))
		assert.Nil(t, job)
	})
}
