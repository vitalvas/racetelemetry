package wire

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalvas/racetelemetry/internal/config"
	"github.com/vitalvas/racetelemetry/internal/model"
)

func TestSink_Name(t *testing.T) {
	s := New(config.SinkEntry{TargetAddr: "127.0.0.1:15000"})
	assert.Equal(t, "wire", s.Name())
}

func TestSink_Run(t *testing.T) {
	skipIfNoNetwork(t)

	t.Run("sends json frames", func(t *testing.T) {
		listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		listener, err := net.ListenUDP("udp", listenAddr)
		require.NoError(t, err)
		defer listener.Close()

		s := New(config.SinkEntry{TargetAddr: listener.LocalAddr().String()})

		in := make(chan model.TelemetryFrame, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(ctx, in)
		}()

		time.Sleep(50 * time.Millisecond)

		in <- model.TelemetryFrame{
			SourceName: "test",
			SourceType: "forza",
			IsRaceOn:   model.Ptr(true),
			EngineRPM:  model.Ptr(float32(5000)),
		}

		buf := make([]byte, 65535)

		err = listener.SetReadDeadline(time.Now().Add(time.Second))
		require.NoError(t, err)

		n, _, err := listener.ReadFromUDP(buf)
		require.NoError(t, err)

		var frame model.TelemetryFrame
		require.NoError(t, json.Unmarshal(buf[:n], &frame))

		assert.Equal(t, "test", frame.SourceName)
		assert.Equal(t, "forza", frame.SourceType)
		assert.True(t, *frame.IsRaceOn)
		assert.Equal(t, float32(5000), *frame.EngineRPM)

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
		close(in)

		err = s.Run(context.Background(), in)
		assert.NoError(t, err)
	})

	t.Run("resolve error", func(t *testing.T) {
		s := New(config.SinkEntry{TargetAddr: "invalid-addr"})

		in := make(chan model.TelemetryFrame)
		close(in)

		err := s.Run(context.Background(), in)
		assert.Error(t, err)
	})

	t.Run("sends multiple frames", func(t *testing.T) {
		listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		listener, err := net.ListenUDP("udp", listenAddr)
		require.NoError(t, err)
		defer listener.Close()

		s := New(config.SinkEntry{TargetAddr: listener.LocalAddr().String()})

		in := make(chan model.TelemetryFrame, 2)
		in <- model.TelemetryFrame{EngineRPM: model.Ptr(float32(3000))}
		in <- model.TelemetryFrame{EngineRPM: model.Ptr(float32(4000))}
		close(in)

		err = s.Run(context.Background(), in)
		assert.NoError(t, err)
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
