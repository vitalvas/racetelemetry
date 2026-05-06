package pcars2shm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitalvas/racetelemetry/internal/model"
)

func TestNew(t *testing.T) {
	_, err := New()
	assert.Error(t, err, "shared memory sink should fail on non-Windows")
}

func TestSink_Name(t *testing.T) {
	s := &Sink{}
	assert.Equal(t, "pcars2_shm", s.Name())
}

func TestSink_Run(t *testing.T) {
	t.Run("stops on channel close", func(t *testing.T) {
		s := &Sink{}

		in := make(chan model.TelemetryFrame)
		close(in)

		err := s.Run(context.Background(), in)
		assert.NoError(t, err)
	})

	t.Run("stops on context cancel", func(t *testing.T) {
		s := &Sink{}

		in := make(chan model.TelemetryFrame)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := s.Run(ctx, in)
		assert.NoError(t, err)
	})

	t.Run("write returns error on stub", func(t *testing.T) {
		s := &Sink{}

		in := make(chan model.TelemetryFrame, 1)
		in <- model.TelemetryFrame{IsRaceOn: true}
		close(in)

		err := s.Run(context.Background(), in)
		assert.Error(t, err)
	})
}

func TestShmWriter_Stub(t *testing.T) {
	var w shmWriter

	err := w.write(&model.TelemetryFrame{})
	assert.Error(t, err)

	w.close()
}
