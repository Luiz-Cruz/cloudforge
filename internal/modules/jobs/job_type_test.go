package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobType_From(t *testing.T) {
	var jt JobType

	t.Run("Valid CleanupFailedJob lowercase", func(t *testing.T) {
		res, err := jt.From("cleanup_failed_environments")
		assert.NoError(t, err)
		assert.Equal(t, CleanupFailedJob, res)
	})

	t.Run("Valid ArchiveStaleJob uppercase", func(t *testing.T) {
		res, err := jt.From("ARCHIVE_STALE_ENVIRONMENTS")
		assert.NoError(t, err)
		assert.Equal(t, ArchiveStaleJob, res)
	})

	t.Run("Unknown JobType", func(t *testing.T) {
		_, err := jt.From("UNKNOWN_JOB")
		assert.Error(t, err)
		assert.Equal(t, ErrUnknownJobType, err)
	})
}
