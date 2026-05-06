package wire

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"

	"github.com/vitalvas/racetelemetry/internal/model"
)

const maxPacketSize = 65535

// Source receives telemetry frames as JSON over UDP.
type Source struct {
	listenAddr string
}

// New creates a new wire UDP source.
func New(listenAddr string) *Source {
	return &Source{
		listenAddr: listenAddr,
	}
}

// Name returns the source identifier.
func (s *Source) Name() string {
	return "wire"
}

// Run listens for JSON-encoded telemetry frames over UDP.
func (s *Source) Run(ctx context.Context, out chan<- model.TelemetryFrame) error {
	defer close(out)

	addr, err := net.ResolveUDPAddr("udp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("wire source: resolve address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("wire source: listen: %w", err)
	}
	defer conn.Close()

	slog.Info("wire source listening", "addr", s.listenAddr)

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	buf := make([]byte, maxPacketSize)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			return fmt.Errorf("wire source: read: %w", err)
		}

		var frame model.TelemetryFrame
		if err := json.Unmarshal(buf[:n], &frame); err != nil {
			slog.Debug("wire source: unmarshal error", "error", err)
			continue
		}

		select {
		case out <- frame:
		case <-ctx.Done():
			return nil
		}
	}
}
