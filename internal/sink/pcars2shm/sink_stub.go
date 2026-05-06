//go:build !windows

package pcars2shm

import (
	"errors"

	"github.com/vitalvas/racetelemetry/internal/model"
)

type shmWriter struct{}

// New creates a new pCars2 shared memory sink.
// On non-Windows platforms, this returns an error.
func New() (*Sink, error) {
	return nil, errors.New("pcars2shm: shared memory sink is only available on Windows")
}

func (w *shmWriter) write(_ *model.TelemetryFrame) error {
	return errors.New("pcars2shm: not supported on this platform")
}

func (w *shmWriter) close() {}
