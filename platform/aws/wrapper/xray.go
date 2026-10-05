package wrapper

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
)

//go:generate mockgen -source=./xray.go -destination=./mock/xray_mock.go -package=mock
type Tracer interface {
	Capture(ctx context.Context, name string, fn func(ctx context.Context) error) error
}

type xrayTracer struct{}

func NewXRayTracer() Tracer {
	return &xrayTracer{}
}

func (x *xrayTracer) Capture(ctx context.Context, name string, fn func(ctx context.Context) error) error {
	if xray.GetSegment(ctx) == nil {
		return fn(ctx)
	}
	return xray.Capture(ctx, name, fn)
}

type noopTracer struct{}

func NewNoopTracer() Tracer {
	return &noopTracer{}
}

func (n *noopTracer) Capture(ctx context.Context, name string, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
