package wire

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"

	"github.com/vitalvas/racetelemetry/internal/model"
)

// Sink sends telemetry frames as JSON over UDP.
type Sink struct {
	targetAddr string
}

// New creates a new wire UDP sink.
func New(targetAddr string) *Sink {
	return &Sink{
		targetAddr: targetAddr,
	}
}

// Name returns the sink identifier.
func (s *Sink) Name() string {
	return "wire"
}

// Run consumes telemetry frames and sends them as JSON UDP packets.
func (s *Sink) Run(ctx context.Context, in <-chan model.TelemetryFrame) error {
	addr, err := net.ResolveUDPAddr("udp", s.targetAddr)
	if err != nil {
		return fmt.Errorf("wire sink: resolve address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("wire sink: dial: %w", err)
	}
	defer conn.Close()

	slog.Info("wire sink connected", "target", s.targetAddr)

	for {
		select {
		case <-ctx.Done():
			return nil
		case frame, ok := <-in:
			if !ok {
				return nil
			}

			data, err := json.Marshal(frame)
			if err != nil {
				slog.Debug("wire sink: marshal error", "error", err)
				continue
			}

			if _, err := conn.Write(data); err != nil {
				slog.Debug("wire sink: write error", "error", err)
			}
		}
	}
}
