//go:build windows

package pcars1shm

import (
	"errors"

	"github.com/vitalvas/racetelemetry/internal/model"
)

type shmWriter struct{}

// New creates a new pCars1 shared memory sink.
func New() (*Sink, error) {
	return nil, errors.New("pcars1shm: shared memory sink is not yet implemented")
}

func (w *shmWriter) write(_ *model.TelemetryFrame) error {
	return errors.New("pcars1shm: not yet implemented")
}

func (w *shmWriter) close() {}
