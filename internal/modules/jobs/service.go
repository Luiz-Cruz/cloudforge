package jobs

import "context"

//go:generate mockgen -source=./service.go -destination=./mock/jobs_mock.go -package=mock
type Job interface {
	RunJob(ctx context.Context) error
}
