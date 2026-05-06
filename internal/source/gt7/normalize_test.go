package gt7

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeGear_GT7(t *testing.T) {
	tests := []struct {
		name  string
		input uint8
		want  int8
	}{
		{"reverse", 0, -1},
		{"first", 1, 1},
		{"second", 2, 2},
		{"sixth", 6, 6},
		{"neutral nibble", 15, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, normalizeGear(tt.input))
		})
	}
}

func TestNormalize_GT7(t *testing.T) {
	t.Run("full packet normalization", func(t *testing.T) {
		pkt := &Packet{
			PositionX:        10.0,
			PositionY:        20.0,
			PositionZ:        30.0,
			RPM:              7000.0,
			CarSpeed:         40.0,
			CurrentFuel:      30.0,
			FuelCapacity:     60.0,
			OilTemp:          95.0,
			WaterTemp:        85.0,
			TireTemp:         [4]float32{80.0, 81.0, 78.0, 79.0},
			Boost:            1.5,
			CurrentLap:       2,
			TotalLaps:        5,
			CurrentPosition:  3,
			BestLapTime:      90500,
			LastLapTime:      91200,
			CurrentLapTimeMS: 45123,
			CurrentGear:      4,
			Throttle:         255,
			Brake:            128,
			Clutch:           0.5,
			SteeringAngle:    0.25,
			StatusFlags:      0x01,
			RPMRevLimiter:    8500,
			GearRatios:       [8]float32{3.5, 2.5, 1.8, 1.3, 1.0, 0.8, 0.0, 0.0},
		}

		frame := normalize(pkt)

		assert.True(t, *frame.IsRaceOn)
		assert.Equal(t, float32(7000.0), *frame.EngineRPM)
		assert.Equal(t, float32(8500.0), *frame.EngineMaxRPM)
		assert.Equal(t, int8(4), *frame.Gear)
		assert.Equal(t, float32(40.0), *frame.Speed)
		assert.InDelta(t, 1.0, *frame.Throttle, 0.01)
		assert.InDelta(t, 0.502, *frame.Brake, 0.01)
		assert.Equal(t, float32(0.5), *frame.Clutch)
		assert.Equal(t, float32(0.25), *frame.Steer)
		assert.InDelta(t, 0.5, *frame.Fuel, 0.01)
		assert.Equal(t, float32(95.0), *frame.OilTemp)
		assert.Equal(t, float32(85.0), *frame.WaterTemp)
		assert.Equal(t, float32(80.0), frame.TireTemp[0])
		assert.Equal(t, float32(0.5), *frame.Boost)
		assert.Equal(t, int16(2), *frame.LapNumber)
		assert.Equal(t, int16(5), *frame.TotalLaps)
		assert.Equal(t, uint8(3), *frame.RacePosition)
		assert.InDelta(t, 90.5, *frame.BestLapTime, 0.01)
		assert.InDelta(t, 91.2, *frame.LastLapTime, 0.01)
		assert.InDelta(t, 45.123, *frame.CurrentLapTime, 0.001)
		assert.Equal(t, [8]float32{3.5, 2.5, 1.8, 1.3, 1.0, 0.8, 0.0, 0.0}, *frame.GearRatios)
	})

	t.Run("zero fuel capacity", func(t *testing.T) {
		pkt := &Packet{
			CurrentFuel:  10.0,
			FuelCapacity: 0,
		}

		frame := normalize(pkt)
		assert.Nil(t, frame.Fuel)
	})

	t.Run("race not on", func(t *testing.T) {
		frame := normalize(&Packet{StatusFlags: 0})
		assert.False(t, *frame.IsRaceOn)
	})

	t.Run("no current lap time in standard format", func(t *testing.T) {
		pkt := &Packet{
			CurrentLapTimeMS: 0,
		}

		frame := normalize(pkt)
		assert.Nil(t, frame.CurrentLapTime)
	})
}

func BenchmarkNormalize_GT7(b *testing.B) {
	pkt := &Packet{
		PositionX:        10.0,
		PositionY:        20.0,
		PositionZ:        30.0,
		RPM:              7000.0,
		CarSpeed:         40.0,
		CurrentFuel:      30.0,
		FuelCapacity:     60.0,
		OilTemp:          95.0,
		WaterTemp:        85.0,
		TireTemp:         [4]float32{80.0, 81.0, 78.0, 79.0},
		Boost:            1.5,
		CurrentLap:       2,
		TotalLaps:        5,
		StatusFlags:      0x01,
		RPMRevLimiter:    8500,
		CurrentGear:      4,
		Throttle:         255,
		Brake:            128,
		Clutch:           0.5,
		SteeringAngle:    0.25,
		CurrentLapTimeMS: 45123,
		GearRatios:       [8]float32{3.5, 2.5, 1.8, 1.3, 1.0, 0.8, 0.0, 0.0},
	}

	b.ResetTimer()

	for b.Loop() {
		normalize(pkt)
	}
}
