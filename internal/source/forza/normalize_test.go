package forza

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeGear(t *testing.T) {
	tests := []struct {
		name  string
		input uint8
		want  int8
	}{
		{"reverse", 0, -1},
		{"neutral", 1, 0},
		{"first", 2, 1},
		{"second", 3, 2},
		{"sixth", 7, 6},
		{"tenth", 11, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, normalizeGear(tt.input))
		})
	}
}

func TestFahrenheitToCelsius(t *testing.T) {
	tests := []struct {
		name       string
		fahrenheit float32
		celsius    float32
	}{
		{"freezing", 32.0, 0.0},
		{"boiling", 212.0, 100.0},
		{"body temp", 98.6, 37.0},
		{"negative", -40.0, -40.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.celsius, fahrenheitToCelsius(tt.fahrenheit), 0.1)
		})
	}
}

func TestNormalize(t *testing.T) {
	t.Run("full packet normalization", func(t *testing.T) {
		pkt := &Packet{
			IsRaceOn:               1,
			CurrentRPM:             5000.0,
			EngineMaxRPM:           8000.0,
			EngineIdleRPM:          900.0,
			Speed:                  30.0,
			Accel:                  255,
			Brake:                  128,
			Clutch:                 0,
			HandBrake:              0,
			Gear:                   3,
			Steer:                  -64,
			TireTemp:               [4]float32{212.0, 212.0, 200.0, 200.0},
			SuspensionTravelMeters: [4]float32{0.1, 0.1, 0.12, 0.12},
			PositionX:              100.0,
			PositionY:              50.0,
			PositionZ:              200.0,
			Fuel:                   0.75,
			LapNumber:              3,
			RacePosition:           5,
			BestLap:                65.5,
			LastLap:                66.2,
			CurrentLap:             30.1,
		}

		frame := normalize(pkt)

		assert.True(t, frame.IsRaceOn)
		assert.Equal(t, float32(5000.0), frame.EngineRPM)
		assert.Equal(t, float32(8000.0), frame.EngineMaxRPM)
		assert.Equal(t, float32(900.0), frame.EngineIdleRPM)
		assert.Equal(t, int8(2), frame.Gear)
		assert.InDelta(t, 1.0, frame.Throttle, 0.01)
		assert.InDelta(t, 0.502, frame.Brake, 0.01)
		assert.InDelta(t, -0.504, frame.Steer, 0.01)
		assert.InDelta(t, 100.0, frame.TireTemp[0], 0.1)
		assert.Equal(t, int16(3), frame.LapNumber)
		assert.Equal(t, uint8(5), frame.RacePosition)
		assert.Equal(t, float32(0.75), frame.Fuel)
		assert.Equal(t, float32(100.0), frame.PositionX)
		assert.Equal(t, float32(50.0), frame.PositionY)
		assert.Equal(t, float32(200.0), frame.PositionZ)
	})

	t.Run("race not on", func(t *testing.T) {
		frame := normalize(&Packet{IsRaceOn: 0})
		assert.False(t, frame.IsRaceOn)
	})
}

func BenchmarkNormalize(b *testing.B) {
	pkt := &Packet{
		IsRaceOn:               1,
		CurrentRPM:             5000.0,
		EngineMaxRPM:           8000.0,
		EngineIdleRPM:          900.0,
		Speed:                  30.0,
		Accel:                  255,
		Brake:                  128,
		Gear:                   3,
		Steer:                  -64,
		TireTemp:               [4]float32{212.0, 212.0, 200.0, 200.0},
		SuspensionTravelMeters: [4]float32{0.1, 0.1, 0.12, 0.12},
		PositionX:              100.0,
		PositionY:              50.0,
		PositionZ:              200.0,
		Fuel:                   0.75,
		LapNumber:              3,
		RacePosition:           5,
	}

	b.ResetTimer()

	for b.Loop() {
		normalize(pkt)
	}
}
