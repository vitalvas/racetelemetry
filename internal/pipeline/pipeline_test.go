package pipeline

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalvas/racetelemetry/internal/model"
)

type mockSource struct {
	name   string
	frames []model.TelemetryFrame
}

func (s *mockSource) Name() string { return s.name }

func (s *mockSource) Run(_ context.Context, out chan<- model.TelemetryFrame) error {
	defer close(out)

	for _, f := range s.frames {
		out <- f
	}

	return nil
}

type mockSink struct {
	mu       sync.Mutex
	received []model.TelemetryFrame
}

func (s *mockSink) Name() string { return "mock_sink" }

func (s *mockSink) Run(_ context.Context, in <-chan model.TelemetryFrame) error {
	for f := range in {
		s.mu.Lock()
		s.received = append(s.received, f)
		s.mu.Unlock()
	}

	return nil
}

func (s *mockSink) getReceived() []model.TelemetryFrame {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]model.TelemetryFrame, len(s.received))
	copy(result, s.received)

	return result
}

func TestPipeline_SingleSourceSingleSink(t *testing.T) {
	source := &mockSource{
		name:   "src1",
		frames: []model.TelemetryFrame{{EngineRPM: model.Ptr(float32(5000))}, {EngineRPM: model.Ptr(float32(6000))}},
	}
	sink := &mockSink{}

	pipe := New(
		map[string]Source{"src1": source},
		map[string]SinkRoute{"sink1": {Sink: sink, Inputs: []string{"src1"}}},
		slog.Default(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := pipe.Run(ctx)
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	received := sink.getReceived()
	assert.Len(t, received, 2)
	assert.Equal(t, float32(5000), *received[0].EngineRPM)
	assert.Equal(t, float32(6000), *received[1].EngineRPM)
}

func TestPipeline_FanOutToMultipleSinks(t *testing.T) {
	source := &mockSource{
		name:   "src1",
		frames: []model.TelemetryFrame{{EngineRPM: model.Ptr(float32(5000))}, {EngineRPM: model.Ptr(float32(6000))}, {EngineRPM: model.Ptr(float32(7000))}},
	}
	sink1 := &mockSink{}
	sink2 := &mockSink{}

	pipe := New(
		map[string]Source{"src1": source},
		map[string]SinkRoute{
			"sink1": {Sink: sink1, Inputs: []string{"src1"}},
			"sink2": {Sink: sink2, Inputs: []string{"src1"}},
		},
		slog.Default(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := pipe.Run(ctx)
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	assert.Len(t, sink1.getReceived(), 3)
	assert.Len(t, sink2.getReceived(), 3)
}

func TestPipeline_RoutedSinks(t *testing.T) {
	srcA := &mockSource{name: "srcA", frames: []model.TelemetryFrame{{EngineRPM: model.Ptr(float32(1000))}}}
	srcB := &mockSource{name: "srcB", frames: []model.TelemetryFrame{{EngineRPM: model.Ptr(float32(2000))}}}

	sinkAll := &mockSink{}
	sinkBOnly := &mockSink{}

	pipe := New(
		map[string]Source{"srcA": srcA, "srcB": srcB},
		map[string]SinkRoute{
			"all":    {Sink: sinkAll, Inputs: []string{"srcA", "srcB"}},
			"b_only": {Sink: sinkBOnly, Inputs: []string{"srcB"}},
		},
		slog.Default(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := pipe.Run(ctx)
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	allReceived := sinkAll.getReceived()
	assert.Len(t, allReceived, 2)

	rpms := []float32{*allReceived[0].EngineRPM, *allReceived[1].EngineRPM}
	assert.Contains(t, rpms, float32(1000))
	assert.Contains(t, rpms, float32(2000))

	bOnlyReceived := sinkBOnly.getReceived()
	assert.Len(t, bOnlyReceived, 1)
	assert.Equal(t, float32(2000), *bOnlyReceived[0].EngineRPM)
}

func TestPipeline_ContextCancel(t *testing.T) {
	source := &mockSource{name: "src1", frames: []model.TelemetryFrame{{IsRaceOn: model.Ptr(true)}}}
	sink := &mockSink{}

	pipe := New(
		map[string]Source{"src1": source},
		map[string]SinkRoute{"sink1": {Sink: sink, Inputs: []string{"src1"}}},
		slog.Default(),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := pipe.Run(ctx)
	assert.NoError(t, err)
}
