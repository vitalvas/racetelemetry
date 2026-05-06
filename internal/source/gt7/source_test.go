package gt7

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

	"golang.org/x/crypto/salsa20"
)

func TestSource_Name(t *testing.T) {
	s := New("192.168.1.1", ":33740")
	assert.Equal(t, "gt7", s.Name())
}

func TestSource_Run(t *testing.T) {
	skipIfNoNetwork(t)
	t.Run("receives and normalizes packet", func(t *testing.T) {
		// Find a free port for listening
		listenConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
		require.NoError(t, err)

		listenAddr := listenConn.LocalAddr().String()
		listenConn.Close()

		// Find a free port for heartbeat target
		heartbeatConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
		require.NoError(t, err)

		heartbeatAddr := heartbeatConn.LocalAddr().(*net.UDPAddr)
		heartbeatConn.Close()

		s := New(heartbeatAddr.IP.String(), listenAddr)
		s.consoleAddr = heartbeatAddr.IP.String()

		// Override heartbeat port for test
		out := make(chan model.TelemetryFrame, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(ctx, out)
		}()

		time.Sleep(100 * time.Millisecond)

		// Build and send an encrypted GT7 packet
		encrypted := buildEncryptedGT7Packet(t)

		udpAddr, err := net.ResolveUDPAddr("udp", listenAddr)
		require.NoError(t, err)

		sendConn, err := net.DialUDP("udp", nil, udpAddr)
		require.NoError(t, err)
		defer sendConn.Close()

		_, err = sendConn.Write(encrypted)
		require.NoError(t, err)

		select {
		case frame := <-out:
			assert.True(t, frame.IsRaceOn)
			assert.Equal(t, float32(6500.0), frame.EngineRPM)
		case <-ctx.Done():
			t.Fatal("timeout waiting for frame")
		}

		cancel()
	})

	t.Run("context cancellation", func(t *testing.T) {
		s := New("127.0.0.1", "127.0.0.1:0")

		out := make(chan model.TelemetryFrame, 1)
		ctx, cancel := context.WithCancel(context.Background())

		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(ctx, out)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		select {
		case err := <-errCh:
			assert.NoError(t, err)
		case <-time.After(2 * time.Second):
			t.Fatal("source did not stop after cancel")
		}
	})

	t.Run("invalid listen address", func(t *testing.T) {
		s := New("127.0.0.1", "invalid-addr")

		out := make(chan model.TelemetryFrame, 1)
		ctx := context.Background()

		err := s.Run(ctx, out)
		assert.Error(t, err)
	})
}

func buildEncryptedGT7Packet(t *testing.T) []byte {
	t.Helper()

	ivValue := uint32(1)
	iv2 := ivValue ^ heartbeatIVSeed

	var nonce [8]byte

	binary.LittleEndian.PutUint32(nonce[0:4], iv2)
	binary.LittleEndian.PutUint32(nonce[4:8], ivValue)

	cleartext := make([]byte, PacketSizeAddendum3)

	binary.LittleEndian.PutUint32(cleartext[0:4], magicNumber)
	binary.LittleEndian.PutUint32(cleartext[0x3C:0x40], floatBitsTest(6500.0))
	binary.LittleEndian.PutUint16(cleartext[0x8E:0x90], 0x01)
	binary.LittleEndian.PutUint32(cleartext[0x70:0x74], 1)

	encrypted := make([]byte, len(cleartext))
	salsa20.XORKeyStream(encrypted, cleartext, nonce[:], getSalsa20Key())

	binary.LittleEndian.PutUint32(encrypted[ivOffset:ivOffset+4], ivValue)

	return encrypted
}

func floatBitsTest(f float32) uint32 {
	return math.Float32bits(f)
}

func skipIfNoNetwork(t *testing.T) {
	t.Helper()

	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("skipping: UDP sockets unavailable: %v", err)
	}

	conn.Close()
}
