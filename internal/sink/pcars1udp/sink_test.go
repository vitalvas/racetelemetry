package pcars1udp

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitalvas/racetelemetry/internal/model"
)

func TestSink_Name(t *testing.T) {
	s := New("127.0.0.1:5606")
	assert.Equal(t, "pcars1_udp", s.Name())
}

func TestSink_Run(t *testing.T) {
	skipIfNoNetwork(t)
	t.Run("sends telemetry packets", func(t *testing.T) {
		listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		listener, err := net.ListenUDP("udp", listenAddr)
		require.NoError(t, err)
		defer listener.Close()

		s := New(listener.LocalAddr().String())

		in := make(chan model.TelemetryFrame, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(ctx, in)
		}()

		time.Sleep(50 * time.Millisecond)

		in <- model.TelemetryFrame{
			IsRaceOn:  model.Ptr(true),
			EngineRPM: model.Ptr(float32(5000.0)),
			Speed:     model.Ptr(float32(30.0)),
			Throttle:  model.Ptr(float32(0.75)),
			Brake:     model.Ptr(float32(0.5)),
			Gear:      model.Ptr(int8(3)),
		}

		buf := make([]byte, 2048)

		err = listener.SetReadDeadline(time.Now().Add(time.Second))
		require.NoError(t, err)

		n, _, err := listener.ReadFromUDP(buf)
		require.NoError(t, err)
		assert.Equal(t, telemetryPacketSize, n)

		assert.Equal(t, uint8(PacketTypeCarPhysics), buf[10])
		assert.Equal(t, uint8(PacketVersion), buf[11])

		cancel()
	})

	t.Run("stops on channel close", func(t *testing.T) {
		listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		require.NoError(t, err)

		listener, err := net.ListenUDP("udp", listenAddr)
		require.NoError(t, err)
		defer listener.Close()

		s := New(listener.LocalAddr().String())

		in := make(chan model.TelemetryFrame)
		errCh := make(chan error, 1)

		go func() {
			errCh <- s.Run(context.Background(), in)
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

		s := New(listener.LocalAddr().String())

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
