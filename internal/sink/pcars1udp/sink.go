package pcars1udp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync/atomic"
	"time"

	"github.com/vitalvas/racetelemetry/internal/config"
	"github.com/vitalvas/racetelemetry/internal/model"
)

const statePacketInterval = time.Second

// Sink sends telemetry as pCars1 UDP packets.
type Sink struct {
	targetAddr string
	seqNum     atomic.Uint32
}

// New creates a new pCars1 UDP sink.
func New(entry config.SinkEntry) *Sink {
	return &Sink{
		targetAddr: entry.TargetAddr,
	}
}

// Name returns the sink identifier.
func (s *Sink) Name() string {
	return "pcars1_udp"
}

// Run consumes telemetry frames and sends pCars1 UDP packets.
func (s *Sink) Run(ctx context.Context, in <-chan model.TelemetryFrame) error {
	addr, err := net.ResolveUDPAddr("udp", s.targetAddr)
	if err != nil {
		return fmt.Errorf("pcars1udp: resolve address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("pcars1udp: dial: %w", err)
	}
	defer conn.Close()

	slog.Info("pcars1 udp sink connected", "target", s.targetAddr)

	var lastStatePacket time.Time

	for {
		select {
		case <-ctx.Done():
			return nil

		case frame, ok := <-in:
			if !ok {
				return nil
			}

			seq := s.seqNum.Add(1)

			telemetryData := convertTelemetry(seq, &frame)
			if _, err := conn.Write(telemetryData); err != nil {
				slog.Debug("pcars1udp: write telemetry", "error", err)
			}

			if time.Since(lastStatePacket) >= statePacketInterval {
				lastStatePacket = time.Now()

				gameStateData := convertGameState(seq, &frame)
				if _, err := conn.Write(gameStateData); err != nil {
					slog.Debug("pcars1udp: write game state", "error", err)
				}

				timingsData := convertTimings(seq, &frame)
				if _, err := conn.Write(timingsData); err != nil {
					slog.Debug("pcars1udp: write timings", "error", err)
				}
			}
		}
	}
}
