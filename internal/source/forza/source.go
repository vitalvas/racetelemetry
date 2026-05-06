package forza

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/vitalvas/racetelemetry/internal/model"
)

// Source reads Forza telemetry via UDP.
// Supports Forza Motorsport 7 (V1), Forza Horizon 4/5 (V2),
// and Forza Motorsport 2023 (V3) by auto-detecting packet size.
type Source struct {
	listenAddr string
}

// New creates a new Forza telemetry source.
func New(listenAddr string) *Source {
	return &Source{
		listenAddr: listenAddr,
	}
}

// Name returns the source identifier.
func (s *Source) Name() string {
	return "forza"
}

// Run listens for Forza UDP telemetry packets and sends normalized frames.
func (s *Source) Run(ctx context.Context, out chan<- model.TelemetryFrame) error {
	defer close(out)

	addr, err := net.ResolveUDPAddr("udp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("forza: resolve address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("forza: listen: %w", err)
	}
	defer conn.Close()

	slog.Info("forza source listening", "addr", s.listenAddr)

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	buf := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			return fmt.Errorf("forza: read: %w", err)
		}

		if n != PacketSizeV1 && n != PacketSizeV2 && n != PacketSizeV3 {
			continue
		}

		pkt, err := decodePacket(buf[:n])
		if err != nil {
			slog.Debug("forza: decode error", "error", err)
			continue
		}

		frame := normalize(pkt)

		select {
		case out <- frame:
		case <-ctx.Done():
			return nil
		}
	}
}
