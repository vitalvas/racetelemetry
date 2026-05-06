package forza

import (
	"context"
	"encoding/binary"
	"math"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalvas/racetelemetry/internal/model"
)

func TestSource_Name(t *testing.T) {
	s := New(":0")
	assert.Equal(t, "forza", s.Name())
}

func TestSource_Run(t *testing.T) {
	skipIfNoNetwork(t)

	t.Run("receives and normalizes packet", func(t *testing.T) {
		// Use port 0 to get a random available port
		s := New("127.0.0.1:0")

		out := make(chan model.TelemetryFrame, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Start source in a goroutine - we need to discover the actual port
		listenAddr := make(chan string, 1)

		// Resolve address to find actual port
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		conn, err := net.ListenUDP("udp", addr)
		require.NoError(t, err)

		actualAddr := conn.LocalAddr().String()
		conn.Close()

		s.listenAddr = actualAddr
		listenAddr <- actualAddr

		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(ctx, out)
		}()

		// Give the source time to start listening
		time.Sleep(50 * time.Millisecond)

		// Send a valid Forza packet
		sendAddr := <-listenAddr
		udpAddr, err := net.ResolveUDPAddr("udp", sendAddr)
		require.NoError(t, err)

		sendConn, err := net.DialUDP("udp", nil, udpAddr)
		require.NoError(t, err)
		defer sendConn.Close()

		pkt := buildTestPacket()
		_, err = sendConn.Write(pkt)
		require.NoError(t, err)

		select {
		case frame := <-out:
			assert.True(t, *frame.IsRaceOn)
			assert.Equal(t, float32(5000.0), *frame.EngineRPM)
			assert.Equal(t, float32(30.0), *frame.Speed)
		case <-ctx.Done():
			t.Fatal("timeout waiting for frame")
		}

		cancel()
	})

	t.Run("ignores wrong size packets", func(t *testing.T) {
		s := New("127.0.0.1:0")

		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		conn, err := net.ListenUDP("udp", addr)
		require.NoError(t, err)

		s.listenAddr = conn.LocalAddr().String()
		conn.Close()

		out := make(chan model.TelemetryFrame, 1)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		go func() {
			_ = s.Run(ctx, out)
		}()

		time.Sleep(50 * time.Millisecond)

		udpAddr, err := net.ResolveUDPAddr("udp", s.listenAddr)
		require.NoError(t, err)

		sendConn, err := net.DialUDP("udp", nil, udpAddr)
		require.NoError(t, err)
		defer sendConn.Close()

		// Send wrong size
		_, err = sendConn.Write([]byte{1, 2, 3})
		require.NoError(t, err)

		select {
		case <-out:
			t.Fatal("should not receive frame for wrong size packet")
		case <-time.After(200 * time.Millisecond):
			// Expected: no frame received
		}

		cancel()
	})
}

func buildTestPacket() []byte {
	buf := make([]byte, PacketSizeV2)

	// IsRaceOn = 1
	binary.LittleEndian.PutUint32(buf[0:4], 1)
	// CurrentRPM = 5000.0
	binary.LittleEndian.PutUint32(buf[16:20], math.Float32bits(5000.0))
	// Speed at V2 offset: posOffset(244) + 12
	binary.LittleEndian.PutUint32(buf[256:260], math.Float32bits(30.0))
	// Gear at V2 offset: posOffset(244) + 12 + 12 + 16 + 28 + 7 = 319
	buf[319] = 3

	return buf
}

func skipIfNoNetwork(t *testing.T) {
	t.Helper()

	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("skipping: UDP sockets unavailable: %v", err)
	}

	conn.Close()
}
