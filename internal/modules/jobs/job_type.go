package jobs

import (
	"errors"
	"strings"
)

type JobType string

const (
	CleanupFailedJob JobType = "CLEANUP_FAILED_ENVIRONMENTS"
	ArchiveStaleJob  JobType = "ARCHIVE_STALE_ENVIRONMENTS"
)

var (
	ErrUnknownJobType = errors.New("unknown job type")
)

func (j *JobType) From(val string) (JobType, error) {
	upper := strings.ToUpper(strings.TrimSpace(val))
	switch JobType(upper) {
	case CleanupFailedJob:
		*j = CleanupFailedJob
		return CleanupFailedJob, nil
	case ArchiveStaleJob:
		*j = ArchiveStaleJob
		return ArchiveStaleJob, nil
	default:
		return "", ErrUnknownJobType
	}
}
