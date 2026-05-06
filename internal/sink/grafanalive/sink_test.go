package grafanalive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalvas/racetelemetry/internal/model"
)

func TestSink_Name(t *testing.T) {
	s := New("http://localhost:3000/api/live/push/test", "test-key")
	assert.Equal(t, "grafana_live", s.Name())
}

func TestSink_Run(t *testing.T) {
	t.Run("pushes telemetry", func(t *testing.T) {
		var receivedBody string
		var receivedAuth string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")

			body, _ := io.ReadAll(r.Body)
			receivedBody = string(body)

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		s := New(fmt.Sprintf("%s/api/live/push/race", server.URL), "my-api-key")

		in := make(chan model.TelemetryFrame, 1)
		in <- model.TelemetryFrame{
			Timestamp: time.Unix(1700000000, 0),
			IsRaceOn:  true,
			EngineRPM: 5000,
			Speed:     30.0,
			Throttle:  0.75,
			Gear:      3,
			TireTemp:  [4]float32{85.0, 86.0, 82.0, 83.0},
		}
		close(in)

		err := s.Run(context.Background(), in)
		require.NoError(t, err)

		assert.Equal(t, "Bearer my-api-key", receivedAuth)
		assert.Contains(t, receivedBody, "telemetry ")
		assert.Contains(t, receivedBody, "engine_rpm=5000.00")
		assert.Contains(t, receivedBody, "speed=30.0000")
		assert.Contains(t, receivedBody, "gear=3i")
		assert.Contains(t, receivedBody, "is_race_on=true")

		// Per-wheel data
		assert.Contains(t, receivedBody, "tire,wheel=fl temp=85.0")
		assert.Contains(t, receivedBody, "tire,wheel=fr temp=86.0")
		assert.Contains(t, receivedBody, "tire,wheel=rl temp=82.0")
		assert.Contains(t, receivedBody, "tire,wheel=rr temp=83.0")
	})

	t.Run("handles server error gracefully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		s := New(fmt.Sprintf("%s/api/live/push/race", server.URL), "key")

		in := make(chan model.TelemetryFrame, 1)
		in <- model.TelemetryFrame{IsRaceOn: true, EngineRPM: 5000}
		close(in)

		// Should not return error - just logs debug
		err := s.Run(context.Background(), in)
		assert.NoError(t, err)
	})

	t.Run("stops on context cancel", func(t *testing.T) {
		s := New("http://localhost:1/noop", "key")

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
		s := New("http://localhost:1/noop", "key")

		in := make(chan model.TelemetryFrame)
		close(in)

		err := s.Run(context.Background(), in)
		assert.NoError(t, err)
	})
}

func TestFormatLineProtocol(t *testing.T) {
	frame := &model.TelemetryFrame{
		Timestamp: time.Unix(1700000000, 0),
		IsRaceOn:  true,
		EngineRPM: 5000,
		Speed:     30.5,
		Gear:      3,
		TireTemp:  [4]float32{85.0, 86.0, 82.0, 83.0},
	}

	line := formatLineProtocol(frame)
	lines := strings.Split(strings.TrimSpace(line), "\n")

	assert.Len(t, lines, 5) // 1 telemetry + 4 tire measurements
	assert.True(t, strings.HasPrefix(lines[0], "telemetry "))
	assert.True(t, strings.HasPrefix(lines[1], "tire,wheel=fl "))
	assert.True(t, strings.HasPrefix(lines[2], "tire,wheel=fr "))
	assert.True(t, strings.HasPrefix(lines[3], "tire,wheel=rl "))
	assert.True(t, strings.HasPrefix(lines[4], "tire,wheel=rr "))
}

func BenchmarkFormatLineProtocol(b *testing.B) {
	frame := &model.TelemetryFrame{
		Timestamp: time.Unix(1700000000, 0),
		IsRaceOn:  true,
		EngineRPM: 5000,
		Speed:     30.0,
		Throttle:  0.75,
		Brake:     0.5,
		Gear:      3,
		TireTemp:  [4]float32{85.0, 86.0, 82.0, 83.0},
	}

	b.ResetTimer()

	for b.Loop() {
		formatLineProtocol(frame)
	}
}
