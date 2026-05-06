package pcars2udp

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalvas/racetelemetry/internal/config"
	"github.com/vitalvas/racetelemetry/internal/model"
)

func TestSink_Name(t *testing.T) {
	s := New(config.SinkEntry{TargetAddr: "127.0.0.1:5606"})
	assert.Equal(t, "pcars2_udp", s.Name())
}

func TestSink_Run(t *testing.T) {
	skipIfNoNetwork(t)
	t.Run("sends telemetry packets", func(t *testing.T) {
		// Listen for UDP packets
		listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		listener, err := net.ListenUDP("udp", listenAddr)
		require.NoError(t, err)
		defer listener.Close()

		targetAddr := listener.LocalAddr().String()
		s := New(config.SinkEntry{TargetAddr: targetAddr})

		in := make(chan model.TelemetryFrame, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(ctx, in)
		}()

		time.Sleep(50 * time.Millisecond)

		// Send a frame
		in <- model.TelemetryFrame{
			IsRaceOn:  model.Ptr(true),
			EngineRPM: model.Ptr(float32(5000.0)),
			Speed:     model.Ptr(float32(30.0)),
			Throttle:  model.Ptr(float32(0.75)),
			Brake:     model.Ptr(float32(0.5)),
			Gear:      model.Ptr(int8(3)),
		}

		// Read the telemetry packet
		buf := make([]byte, 2048)

		err = listener.SetReadDeadline(time.Now().Add(time.Second))
		require.NoError(t, err)

		n, _, err := listener.ReadFromUDP(buf)
		require.NoError(t, err)
		assert.Equal(t, telemetryPacketSize, n)

		// Verify header
		assert.Equal(t, uint8(PacketTypeCarPhysics), buf[10])

		cancel()
	})

	t.Run("stops on channel close", func(t *testing.T) {
		listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		listener, err := net.ListenUDP("udp", listenAddr)
		require.NoError(t, err)
		defer listener.Close()

		s := New(config.SinkEntry{TargetAddr: listener.LocalAddr().String()})

		in := make(chan model.TelemetryFrame)
		ctx := context.Background()

		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(ctx, in)
		}()

		time.Sleep(50 * time.Millisecond)
		close(in)

		select {
		case err := <-errCh:
			assert.NoError(t, err)
		case <-time.After(2 * time.Second):
			t.Fatal("sink did not stop after channel close")
		}
	})

	t.Run("stops on context cancel", func(t *testing.T) {
		listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		listener, err := net.ListenUDP("udp", listenAddr)
		require.NoError(t, err)
		defer listener.Close()

		s := New(config.SinkEntry{TargetAddr: listener.LocalAddr().String()})

		in := make(chan model.TelemetryFrame)
		ctx, cancel := context.WithCancel(context.Background())

		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(ctx, in)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		select {
		case err := <-errCh:
			assert.NoError(t, err)
		case <-time.After(2 * time.Second):
			t.Fatal("sink did not stop after cancel")
		}
	})
}

func skipIfNoNetwork(t *testing.T) {
	t.Helper()

	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("skipping: UDP sockets unavailable: %v", err)
	}

	conn.Close()
}
