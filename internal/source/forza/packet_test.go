package forza

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodePacket(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		posOffset int // where PositionX starts in the wire format
	}{
		{"V1 forza motorsport 7", PacketSizeV1, 232},
		{"V2 forza horizon 4/5", PacketSizeV2, 244},
		{"V3 forza motorsport 2023", PacketSizeV3, 232},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, tt.size)

			// Common fields (same offset across all versions)
			binary.LittleEndian.PutUint32(buf[0:4], 1) // IsRaceOn
			binary.LittleEndian.PutUint32(buf[4:8], 12345)
			putTestFloat32(buf, 8, 8000.0)  // EngineMaxRPM
			putTestFloat32(buf, 12, 900.0)  // EngineIdleRPM
			putTestFloat32(buf, 16, 5000.0) // CurrentRPM

			// Version-dependent fields
			putTestFloat32(buf, tt.posOffset, 100.0)   // PositionX
			putTestFloat32(buf, tt.posOffset+12, 44.4) // Speed

			// Gear, Accel, Brake, Steer at version-dependent tail
			tailOff := tt.posOffset + 12 + 12 + 16 + 28 // pos+speed/power/torque+tiretemps+boost..racetime
			buf[tailOff+7] = 3                          // Gear
			buf[tailOff+3] = 200                        // Accel
			buf[tailOff+4] = 100                        // Brake
			buf[tailOff+8] = uint8(192)                 // Steer = int8(-64)

			pkt, err := decodePacket(buf)
			require.NoError(t, err)

			assert.Equal(t, int32(1), pkt.IsRaceOn)
			assert.Equal(t, uint32(12345), pkt.TimestampMS)
			assert.Equal(t, float32(8000.0), pkt.EngineMaxRPM)
			assert.Equal(t, float32(900.0), pkt.EngineIdleRPM)
			assert.Equal(t, float32(5000.0), pkt.CurrentRPM)
			assert.Equal(t, float32(100.0), pkt.PositionX)
			assert.Equal(t, float32(44.4), pkt.Speed)
			assert.Equal(t, uint8(3), pkt.Gear)
			assert.Equal(t, uint8(200), pkt.Accel)
			assert.Equal(t, uint8(100), pkt.Brake)
			assert.Equal(t, int8(-64), pkt.Steer)
		})
	}

	t.Run("wrong size", func(t *testing.T) {
		_, err := decodePacket(make([]byte, 100))
		assert.Error(t, err)
	})
}

func BenchmarkDecodePacket(b *testing.B) {
	versions := []struct {
		name string
		size int
	}{
		{"V1", PacketSizeV1},
		{"V2", PacketSizeV2},
		{"V3", PacketSizeV3},
	}

	for _, v := range versions {
		b.Run(v.name, func(b *testing.B) {
			buf := make([]byte, v.size)
			binary.LittleEndian.PutUint32(buf[0:4], 1)
			putTestFloat32(buf, 16, 5000.0)

			b.ResetTimer()

			for b.Loop() {
				_, _ = decodePacket(buf)
			}
		})
	}
}

func FuzzDecodePacket(f *testing.F) {
	f.Add(make([]byte, PacketSizeV1))
	f.Add(make([]byte, PacketSizeV2))
	f.Add(make([]byte, PacketSizeV3))

	seed := make([]byte, PacketSizeV2)
	binary.LittleEndian.PutUint32(seed[0:4], 1)
	putTestFloat32(seed, 16, 5000.0)
	f.Add(seed)

	f.Fuzz(func(t *testing.T, data []byte) {
		pkt, err := decodePacket(data)
		if err != nil {
			return
		}

		if pkt == nil {
			t.Fatal("decodePacket returned nil packet without error")
		}
	})
}

func putTestFloat32(buf []byte, offset int, v float32) {
	binary.LittleEndian.PutUint32(buf[offset:offset+4], math.Float32bits(v))
}
