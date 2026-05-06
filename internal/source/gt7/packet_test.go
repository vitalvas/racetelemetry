package gt7

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodePacket_GT7(t *testing.T) {
	t.Run("standard format", func(t *testing.T) {
		data := make([]byte, PacketSizeStandard)

		putFloat(data, 0x04, 1.0)    // PositionX
		putFloat(data, 0x08, 2.0)    // PositionY
		putFloat(data, 0x0C, 3.0)    // PositionZ
		putFloat(data, 0x28, 0.75)   // Heading
		putFloat(data, 0x3C, 6500.0) // RPM
		putFloat(data, 0x4C, 50.0)   // Speed
		putFloat(data, 0x94, 0.1)    // RoadPlaneX
		putFloat(data, 0x100, 1.5)   // TopSpeedRatio

		putFloat(data, 0x60, 85.0)
		putFloat(data, 0x64, 86.0)
		putFloat(data, 0x68, 82.0)
		putFloat(data, 0x6C, 83.0)

		binary.LittleEndian.PutUint32(data[0x70:0x74], 100)
		binary.LittleEndian.PutUint16(data[0x74:0x76], 3)
		binary.LittleEndian.PutUint16(data[0x76:0x78], 10)

		data[0x90] = 3 | (4 << 4)
		data[0x91] = 200
		data[0x92] = 100

		pkt, err := decodePacket(data)
		require.NoError(t, err)

		assert.Equal(t, float32(1.0), pkt.PositionX)
		assert.Equal(t, float32(2.0), pkt.PositionY)
		assert.Equal(t, float32(3.0), pkt.PositionZ)
		assert.Equal(t, float32(0.75), pkt.Heading)
		assert.Equal(t, float32(6500.0), pkt.RPM)
		assert.Equal(t, float32(50.0), pkt.CarSpeed)
		assert.Equal(t, float32(0.1), pkt.RoadPlaneX)
		assert.Equal(t, float32(1.5), pkt.TopSpeedRatio)
		assert.Equal(t, float32(85.0), pkt.TireTemp[0])
		assert.Equal(t, int32(100), pkt.PackageID)
		assert.Equal(t, int16(3), pkt.CurrentLap)
		assert.Equal(t, int16(10), pkt.TotalLaps)
		assert.Equal(t, uint8(3), pkt.CurrentGear)
		assert.Equal(t, uint8(4), pkt.SuggestedGear)

		// Addendum fields should be zero for standard format
		assert.Equal(t, float32(0), pkt.SteeringAngle)
		assert.Equal(t, int32(0), pkt.CurrentLapTimeMS)
	})

	t.Run("addendum3 format", func(t *testing.T) {
		data := make([]byte, PacketSizeAddendum3)

		putFloat(data, 0x3C, 7000.0)                                    // RPM
		putFloat(data, 0x128, 0.25)                                     // SteeringAngle
		putFloat(data, 0x12C, 0.5)                                      // SteeringVelocity
		binary.LittleEndian.PutUint32(data[0x15C:0x160], uint32(45123)) // CurrentLapTimeMS
		binary.LittleEndian.PutUint32(data[0x70:0x74], 1)

		pkt, err := decodePacket(data)
		require.NoError(t, err)

		assert.Equal(t, float32(7000.0), pkt.RPM)
		assert.Equal(t, float32(0.25), pkt.SteeringAngle)
		assert.Equal(t, float32(0.5), pkt.SteeringVelocity)
		assert.Equal(t, int32(45123), pkt.CurrentLapTimeMS)
	})

	t.Run("packet too short", func(t *testing.T) {
		_, err := decodePacket(make([]byte, 10))
		assert.Error(t, err)
	})
}

func BenchmarkDecodePacket_GT7(b *testing.B) {
	data := make([]byte, PacketSizeAddendum3)
	putFloat(data, 0x3C, 6500.0)
	putFloat(data, 0x4C, 50.0)
	binary.LittleEndian.PutUint32(data[0x70:0x74], 100)
	data[0x90] = 3 | (4 << 4)
	data[0x91] = 200
	data[0x92] = 100

	b.ResetTimer()

	for b.Loop() {
		_, _ = decodePacket(data)
	}
}

func FuzzDecodePacket_GT7(f *testing.F) {
	f.Add(make([]byte, PacketSizeStandard))
	f.Add(make([]byte, PacketSizeAddendum3))

	seed := make([]byte, PacketSizeAddendum3)
	putFloat(seed, 0x3C, 6500.0)
	putFloat(seed, 0x4C, 50.0)
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

func putFloat(buf []byte, offset int, v float32) {
	binary.LittleEndian.PutUint32(buf[offset:offset+4], math.Float32bits(v))
}
