package pipeline

import (
	"context"

	"github.com/vitalvas/racetelemetry/internal/model"
)

// Source reads telemetry from an external game and sends normalized frames.
type Source interface {
	// Run starts receiving telemetry and sends frames to the output channel.
	// It blocks until ctx is cancelled or a fatal error occurs.
	// The source is responsible for closing the output channel when done.
	Run(ctx context.Context, out chan<- model.TelemetryFrame) error

	// Name returns the source identifier for logging.
	Name() string
}
