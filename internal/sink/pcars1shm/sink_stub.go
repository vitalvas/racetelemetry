//go:build !windows

package pcars1shm

import (
	"errors"

	"github.com/vitalvas/racetelemetry/internal/model"
)

type shmWriter struct{}

// New creates a new pCars1 shared memory sink.
// On non-Windows platforms, this returns an error.
func New() (*Sink, error) {
	return nil, errors.New("pcars1shm: shared memory sink is only available on Windows")
}

func (w *shmWriter) write(_ *model.TelemetryFrame) error {
	return errors.New("pcars1shm: not supported on this platform")
}

func (w *shmWriter) close() {}
