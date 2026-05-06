package jsonstdout

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/vitalvas/racetelemetry/internal/model"
)

// Sink writes telemetry frames as JSON lines to stdout.
type Sink struct {
	encoder *json.Encoder
}

// New creates a new JSON stdout sink.
func New() *Sink {
	return newSink(os.Stdout)
}

func newSink(w io.Writer) *Sink {
	return &Sink{
		encoder: json.NewEncoder(w),
	}
}

// Name returns the sink identifier.
func (s *Sink) Name() string {
	return "json_stdout"
}

// Run consumes telemetry frames and writes them as JSON lines to stdout.
func (s *Sink) Run(ctx context.Context, in <-chan model.TelemetryFrame) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case frame, ok := <-in:
			if !ok {
				return nil
			}

			if err := s.encoder.Encode(frame); err != nil {
				return err
			}
		}
	}
}
