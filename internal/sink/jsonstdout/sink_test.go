package jsonstdout

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalvas/racetelemetry/internal/model"
)

func TestSink_Name(t *testing.T) {
	s := New()
	assert.Equal(t, "json_stdout", s.Name())
}

func TestSink_Run(t *testing.T) {
	t.Run("writes json lines", func(t *testing.T) {
		var buf bytes.Buffer
		s := newSink(&buf)

		in := make(chan model.TelemetryFrame, 2)
		in <- model.TelemetryFrame{IsRaceOn: true, EngineRPM: 5000, Speed: 30.0}
		in <- model.TelemetryFrame{IsRaceOn: true, EngineRPM: 6000, Speed: 40.0}
		close(in)

		err := s.Run(context.Background(), in)
		require.NoError(t, err)

		lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
		assert.Len(t, lines, 2)

		var frame1 model.TelemetryFrame
		require.NoError(t, json.Unmarshal(lines[0], &frame1))
		assert.Equal(t, float32(5000), frame1.EngineRPM)
		assert.Equal(t, float32(30.0), frame1.Speed)
		assert.True(t, frame1.IsRaceOn)

		var frame2 model.TelemetryFrame
		require.NoError(t, json.Unmarshal(lines[1], &frame2))
		assert.Equal(t, float32(6000), frame2.EngineRPM)
	})

	t.Run("stops on context cancel", func(t *testing.T) {
		s := newSink(&bytes.Buffer{})

		in := make(chan model.TelemetryFrame)
		ctx, cancel := context.WithCancel(context.Background())

		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(ctx, in)
		}()

		cancel()

		select {
		case err := <-errCh:
			assert.NoError(t, err)
		case <-time.After(time.Second):
			t.Fatal("sink did not stop after cancel")
		}
	})

	t.Run("stops on channel close", func(t *testing.T) {
		s := newSink(&bytes.Buffer{})

		in := make(chan model.TelemetryFrame)
		close(in)

		err := s.Run(context.Background(), in)
		assert.NoError(t, err)
	})
}

func BenchmarkSink_Run(b *testing.B) {
	frame := model.TelemetryFrame{
		IsRaceOn:  true,
		EngineRPM: 5000,
		Speed:     30.0,
		Throttle:  0.75,
		Brake:     0.5,
		Gear:      3,
		TireTemp:  [4]float32{85.0, 86.0, 82.0, 83.0},
	}

	var buf bytes.Buffer
	s := newSink(&buf)

	b.ResetTimer()

	for b.Loop() {
		buf.Reset()
		_ = s.encoder.Encode(frame)
	}
}
