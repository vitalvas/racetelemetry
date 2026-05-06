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

func TestSource_Name(t *testing.T) {
	s := New(config.SourceEntry{ListenAddr: ":15000"})
	assert.Equal(t, "wire", s.Name())
}

func TestSource_Run(t *testing.T) {
	skipIfNoNetwork(t)

	t.Run("receives json frames", func(t *testing.T) {
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		conn, err := net.ListenUDP("udp", addr)
		require.NoError(t, err)

		listenAddr := conn.LocalAddr().String()
		conn.Close()

		s := New(config.SourceEntry{ListenAddr: listenAddr})

		out := make(chan model.TelemetryFrame, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		go func() {
			_ = s.Run(ctx, out)
		}()

		time.Sleep(50 * time.Millisecond)

		// Send a JSON frame
		frame := model.TelemetryFrame{
			SourceName: "upstream",
			SourceType: "forza",
			IsRaceOn:   model.Ptr(true),
			EngineRPM:  model.Ptr(float32(6000)),
			Speed:      model.Ptr(float32(40.0)),
		}

		data, err := json.Marshal(frame)
		require.NoError(t, err)

		udpAddr, err := net.ResolveUDPAddr("udp", listenAddr)
		require.NoError(t, err)

		sendConn, err := net.DialUDP("udp", nil, udpAddr)
		require.NoError(t, err)
		defer sendConn.Close()

		_, err = sendConn.Write(data)
		require.NoError(t, err)

		select {
		case received := <-out:
			assert.Equal(t, "upstream", received.SourceName)
			assert.Equal(t, "forza", received.SourceType)
			assert.True(t, *received.IsRaceOn)
			assert.Equal(t, float32(6000), *received.EngineRPM)
			assert.Equal(t, float32(40.0), *received.Speed)
		case <-ctx.Done():
			t.Fatal("timeout waiting for frame")
		}

		cancel()
	})

	t.Run("skips invalid json", func(t *testing.T) {
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		conn, err := net.ListenUDP("udp", addr)
		require.NoError(t, err)

		listenAddr := conn.LocalAddr().String()
		conn.Close()

		s := New(config.SourceEntry{ListenAddr: listenAddr})

		out := make(chan model.TelemetryFrame, 1)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		go func() {
			_ = s.Run(ctx, out)
		}()

		time.Sleep(50 * time.Millisecond)

		udpAddr, err := net.ResolveUDPAddr("udp", listenAddr)
		require.NoError(t, err)

		sendConn, err := net.DialUDP("udp", nil, udpAddr)
		require.NoError(t, err)
		defer sendConn.Close()

		// Send invalid JSON
		_, err = sendConn.Write([]byte("not json"))
		require.NoError(t, err)

		// Then send valid frame
		frame := model.TelemetryFrame{EngineRPM: model.Ptr(float32(7000))}

		data, err := json.Marshal(frame)
		require.NoError(t, err)

		_, err = sendConn.Write(data)
		require.NoError(t, err)

		select {
		case received := <-out:
			assert.Equal(t, float32(7000), *received.EngineRPM)
		case <-ctx.Done():
			t.Fatal("timeout waiting for frame")
		}

		cancel()
	})

	t.Run("invalid listen address", func(t *testing.T) {
		s := New(config.SourceEntry{ListenAddr: "invalid-addr"})

		out := make(chan model.TelemetryFrame, 1)

		err := s.Run(context.Background(), out)
		assert.Error(t, err)
	})
}

func BenchmarkUnmarshalFrame(b *testing.B) {
	frame := model.TelemetryFrame{
		SourceName: "bench",
		SourceType: "forza",
		IsRaceOn:   model.Ptr(true),
		EngineRPM:  model.Ptr(float32(5000)),
		Speed:      model.Ptr(float32(30.0)),
		TireTemp:   &[4]float32{85.0, 86.0, 82.0, 83.0},
	}

	data, _ := json.Marshal(frame)

	b.ResetTimer()

	for b.Loop() {
		var f model.TelemetryFrame
		_ = json.Unmarshal(data, &f)
	}
}

func skipIfNoNetwork(t *testing.T) {
	t.Helper()

	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("skipping: UDP sockets unavailable: %v", err)
	}

	conn.Close()
}
