package gt7

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/vitalvas/racetelemetry/internal/model"
)

const (
	heartbeatPort     = "33739"
	heartbeatInterval = 10 * time.Second
	heartbeatMessage  = "C" // Request richest format (Addendum3, 368 bytes)
	heartbeatIVSeed   = IVSeedFormatC
	readBufferSize    = 4096
)

// Source reads Gran Turismo telemetry via UDP.
// Supports GT7, GT Sport, and GT6.
type Source struct {
	consoleAddr string
	listenAddr  string
}

// New creates a new GT telemetry source.
func New(consoleAddr, listenAddr string) *Source {
	return &Source{
		consoleAddr: consoleAddr,
		listenAddr:  listenAddr,
	}
}

// Name returns the source identifier.
func (s *Source) Name() string {
	return "gt7"
}

// Run listens for GT UDP telemetry packets and sends normalized frames.
func (s *Source) Run(ctx context.Context, out chan<- model.TelemetryFrame) error {
	defer close(out)

	addr, err := net.ResolveUDPAddr("udp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("gt7: resolve listen address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("gt7: listen: %w", err)
	}
	defer conn.Close()

	heartbeatAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(s.consoleAddr, heartbeatPort))
	if err != nil {
		return fmt.Errorf("gt7: resolve console address: %w", err)
	}

	slog.Info("gt7 source listening", "addr", s.listenAddr, "console", s.consoleAddr)

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	go s.sendHeartbeats(ctx, conn, heartbeatAddr)

	if _, err := conn.WriteToUDP([]byte(heartbeatMessage), heartbeatAddr); err != nil {
		slog.Warn("gt7: initial heartbeat failed", "error", err)
	}

	buf := make([]byte, readBufferSize)
	var lastPackageID int32

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			return fmt.Errorf("gt7: read: %w", err)
		}

		decrypted, err := decrypt(buf[:n], heartbeatIVSeed)
		if err != nil {
			slog.Debug("gt7: decrypt error", "error", err)
			continue
		}

		pkt, err := decodePacket(decrypted)
		if err != nil {
			slog.Debug("gt7: decode error", "error", err)
			continue
		}

		if pkt.PackageID <= lastPackageID {
			continue
		}

		lastPackageID = pkt.PackageID

		frame := normalize(pkt)

		select {
		case out <- frame:
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *Source) sendHeartbeats(ctx context.Context, conn *net.UDPConn, addr *net.UDPAddr) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if _, err := conn.WriteToUDP([]byte(heartbeatMessage), addr); err != nil {
				slog.Debug("gt7: heartbeat failed", "error", err)
			}
		}
	}
}
