package pipeline

import (
	"context"
	"log/slog"
	"sync"

	"github.com/vitalvas/racetelemetry/internal/model"
)

const defaultBufferSize = 16

// SinkRoute describes a sink and the sources it subscribes to.
type SinkRoute struct {
	Sink   Sink
	Inputs []string
}

// Pipeline connects named sources to routed sinks via channels.
type Pipeline struct {
	sources map[string]Source
	sinks   map[string]SinkRoute
	logger  *slog.Logger
	bufSize int
}

// New creates a pipeline with named sources and routed sinks.
func New(sources map[string]Source, sinks map[string]SinkRoute, logger *slog.Logger) *Pipeline {
	return &Pipeline{
		sources: sources,
		sinks:   sinks,
		logger:  logger,
		bufSize: defaultBufferSize,
	}
}

// Run starts all sources and sinks with input routing and blocks
// until completion or cancellation.
func (p *Pipeline) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, len(p.sources)+len(p.sinks))

	// Create per-sink channels
	sinkChs := make(map[string]chan model.TelemetryFrame, len(p.sinks))
	for name := range p.sinks {
		sinkChs[name] = make(chan model.TelemetryFrame, p.bufSize)
	}

	// Build routing table: source name -> list of sink channels to forward to
	routingTable := make(map[string][]chan model.TelemetryFrame)

	for sinkName, route := range p.sinks {
		for _, srcName := range route.Inputs {
			routingTable[srcName] = append(routingTable[srcName], sinkChs[sinkName])
		}
	}

	// Start sources with forwarding to routed sinks
	var fwdWg sync.WaitGroup

	for name, src := range p.sources {
		fwdWg.Add(1)

		srcCh := make(chan model.TelemetryFrame, p.bufSize)
		targets := routingTable[name]

		go func() {
			p.logger.Info("starting source", "name", name)

			if err := src.Run(ctx, srcCh); err != nil {
				errCh <- err
				cancel()
			}
		}()

		go func() {
			defer fwdWg.Done()

			for frame := range srcCh {
				frame.SourceName = name
				frame.SourceType = src.Name()

				for _, ch := range targets {
					trySend(ctx, ch, frame)
				}
			}
		}()
	}

	// Close all sink channels after all forwarders finish
	go func() {
		fwdWg.Wait()

		for _, ch := range sinkChs {
			close(ch)
		}
	}()

	// Start sinks
	var sinkWg sync.WaitGroup

	for name, route := range p.sinks {
		sinkWg.Add(1)

		go func() {
			defer sinkWg.Done()

			p.logger.Info("starting sink", "name", name)

			if err := route.Sink.Run(ctx, sinkChs[name]); err != nil {
				errCh <- err
				cancel()
			}
		}()
	}

	sinkWg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

// trySend attempts a non-blocking send. If the channel is full,
// it drops the oldest frame and sends the new one. For real-time
// telemetry, fresh data is always more valuable than stale data.
func trySend(ctx context.Context, ch chan model.TelemetryFrame, frame model.TelemetryFrame) {
	select {
	case ch <- frame:
		return
	case <-ctx.Done():
		return
	default:
	}

	// Channel full: drop oldest, send new
	select {
	case <-ch:
	default:
	}

	select {
	case ch <- frame:
	case <-ctx.Done():
	}
}
