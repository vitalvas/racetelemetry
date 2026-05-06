package pipeline

import (
	"context"

	"github.com/vitalvas/racetelemetry/internal/model"
)

// Sink consumes telemetry frames and outputs them in a specific format.
type Sink interface {
	// Run starts consuming telemetry frames from the input channel.
	// It blocks until ctx is cancelled, the input channel is closed,
	// or a fatal error occurs.
	Run(ctx context.Context, in <-chan model.TelemetryFrame) error

	// Name returns the sink identifier for logging.
	Name() string
}
