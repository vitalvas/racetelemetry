package pcars1shm

import (
	"context"

	"github.com/vitalvas/racetelemetry/internal/model"
)

// Sink writes telemetry to a pCars1 shared memory region.
type Sink struct {
	writer shmWriter
}

// Name returns the sink identifier.
func (s *Sink) Name() string {
	return "pcars1_shm"
}

// Run consumes telemetry frames and writes to shared memory.
func (s *Sink) Run(ctx context.Context, in <-chan model.TelemetryFrame) error {
	defer s.writer.close()

	for {
		select {
		case <-ctx.Done():
			return nil
		case frame, ok := <-in:
			if !ok {
				return nil
			}

			if err := s.writer.write(&frame); err != nil {
				return err
			}
		}
	}
}
