package pcars1udp

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitalvas/racetelemetry/internal/model"
)

func TestConvertTelemetry(t *testing.T) {
	t.Run("packet size and header", func(t *testing.T) {
		frame := model.TelemetryFrame{
			IsRaceOn:  true,
			EngineRPM: 5000.0,
			Speed:     30.0,
			Throttle:  0.75,
			Brake:     0.5,
			Gear:      3,
		}

		data := convertTelemetry(42, &frame)

		assert.Equal(t, telemetryPacketSize, len(data))

		packetNum := binary.LittleEndian.Uint32(data[0:4])
		assert.Equal(t, uint32(42), packetNum)

		assert.Equal(t, uint8(PacketTypeCarPhysics), data[10])
		assert.Equal(t, uint8(PacketVersion), data[11])
	})

	t.Run("throttle brake encoding", func(t *testing.T) {
		frame := model.TelemetryFrame{
			Throttle: 1.0,
			Brake:    0.0,
		}

		data := convertTelemetry(1, &frame)

		assert.Equal(t, uint8(255), data[13])
		assert.Equal(t, uint8(0), data[14])
	})
}

func TestConvertGameState(t *testing.T) {
	t.Run("race on", func(t *testing.T) {
		frame := model.TelemetryFrame{IsRaceOn: true}
		data := convertGameState(1, &frame)

		assert.Equal(t, gameStatePacketSize, len(data))
		assert.Equal(t, uint8(PacketTypeGameState), data[10])
		assert.Equal(t, uint8(PacketVersion), data[11])

		gameState := binary.LittleEndian.Uint16(data[headerSize : headerSize+2])
		assert.Equal(t, uint16(2), gameState)
	})

	t.Run("race off", func(t *testing.T) {
		frame := model.TelemetryFrame{IsRaceOn: false}
		data := convertGameState(1, &frame)

		gameState := binary.LittleEndian.Uint16(data[headerSize : headerSize+2])
		assert.Equal(t, uint16(0), gameState)
	})
}

func TestConvertTimings(t *testing.T) {
	t.Run("timing data", func(t *testing.T) {
		frame := model.TelemetryFrame{
			IsRaceOn:       true,
			PositionX:      100.0,
			PositionY:      50.0,
			PositionZ:      200.0,
			LapDistance:    1500.0,
			RacePosition:   3,
			LapNumber:      5,
			BestLapTime:    65.5,
			LastLapTime:    66.2,
			CurrentLapTime: 30.1,
		}

		data := convertTimings(1, &frame)

		assert.Equal(t, timingsPacketSize, len(data))
		assert.Equal(t, uint8(PacketTypeTimings), data[10])
		assert.Equal(t, uint8(PacketVersion), data[11])

		assert.Equal(t, uint8(1), data[headerSize])

		posXOffset := headerSize + 1 + 4 + 4 + 4 + 4
		posX := math.Float32frombits(binary.LittleEndian.Uint32(data[posXOffset : posXOffset+4]))
		assert.Equal(t, float32(100.0), posX)
	})
}

func BenchmarkConvertTelemetry(b *testing.B) {
	frame := &model.TelemetryFrame{
		IsRaceOn:     true,
		EngineRPM:    5000.0,
		EngineMaxRPM: 8000.0,
		Speed:        30.0,
		Throttle:     0.75,
		Brake:        0.5,
		Gear:         3,
		TireTemp:     [4]float32{85.0, 86.0, 82.0, 83.0},
		Fuel:         0.75,
	}

	b.ResetTimer()

	for b.Loop() {
		convertTelemetry(1, frame)
	}
}
