package forza

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		assert.True(t, *frame.IsRaceOn)
		assert.Equal(t, float32(5000.0), *frame.EngineRPM)
		assert.Equal(t, float32(8000.0), *frame.EngineMaxRPM)
		assert.Equal(t, float32(900.0), *frame.EngineIdleRPM)
		assert.Equal(t, int8(2), *frame.Gear)
		assert.InDelta(t, 1.0, *frame.Throttle, 0.01)
		assert.InDelta(t, 0.502, *frame.Brake, 0.01)
		assert.InDelta(t, -0.504, *frame.Steer, 0.01)
		require.NotNil(t, frame.TireTemp)
		assert.InDelta(t, 100.0, frame.TireTemp[0], 0.1)
		assert.Equal(t, int16(3), *frame.LapNumber)
		assert.Equal(t, uint8(5), *frame.RacePosition)
		assert.Equal(t, float32(0.75), *frame.Fuel)
		assert.Equal(t, float32(100.0), *frame.PositionX)
		assert.Equal(t, float32(50.0), *frame.PositionY)
		assert.Equal(t, float32(200.0), *frame.PositionZ)
	})

	t.Run("race not on", func(t *testing.T) {
		frame := normalize(&Packet{IsRaceOn: 0})
		assert.False(t, *frame.IsRaceOn)
	})

	// Real Forza Horizon 5 telemetry sample - car idle in neutral
	t.Run("real FH5 idle sample", func(t *testing.T) {
		pkt := &Packet{
			IsRaceOn:         1,
			CurrentRPM:       1120.6162,
			EngineMaxRPM:     7999.995,
			EngineIdleRPM:    1120.6044,
			Speed:            0.0004093854,
			Accel:            0,
			Brake:            0,
			Clutch:           0,
			HandBrake:        0,
			Gear:             1, // Neutral in Forza encoding
			Steer:            0,
			AccelerationX:    0.00055534753,
			AccelerationY:    -0.0066558914,
			AccelerationZ:    -0.0010788094,
			VelocityX:        0.00037223985,
			VelocityY:        0.00005263859,
			VelocityZ:        -0.00016226669,
			AngularVelocityX: 0.000088340166,
			AngularVelocityY: -0.000003308337,
			AngularVelocityZ: -0.00012022164,
			Yaw:              1.8943167,
			Pitch:            0.0639811,
			Roll:             -0.007876748,
			PositionX:        -4.641847,
			PositionY:        369.80182,
			PositionZ:        4243.9,
			Power:            -0.07316666,
			Torque:           -0.0006234875,
			Boost:            -11.024998,
			Fuel:             1.0,
			// TireTemp in Fahrenheit on the wire, celsius after normalize
			TireTemp:               [4]float32{153.22485, 152.06763, 153.59266, 153.59266},
			SuspensionTravelMeters: [4]float32{-0.006121576, -0.0032083988, -0.007178575, -0.00879845},
			WheelRotationSpeed:     [4]float32{0, 0, 0, 0},
			TireSlipRatio:          [4]float32{-0.17556109, -0.07943692, 0.21793158, -0.0019374447},
			TireSlipAngle:          [4]float32{0.12178061, -0.14612146, 0.21287914, -0.20355079},
			CarClass:               2,
			CarOrdinal:             3702,
			LapNumber:              0,
			RacePosition:           0,
		}

		frame := normalize(pkt)

		assert.True(t, *frame.IsRaceOn)
		assert.InDelta(t, 1120.6162, *frame.EngineRPM, 0.01)
		assert.InDelta(t, 7999.995, *frame.EngineMaxRPM, 0.01)
		assert.InDelta(t, 1120.6044, *frame.EngineIdleRPM, 0.01)
		assert.Equal(t, int8(0), *frame.Gear) // Neutral
		assert.InDelta(t, 0.0004093854, *frame.Speed, 0.0001)
		assert.InDelta(t, 0.0, *frame.Throttle, 0.01)
		assert.InDelta(t, 0.0, *frame.Brake, 0.01)
		assert.InDelta(t, 0.0, *frame.Steer, 0.01)

		// Verify F->C conversion: 153.22485F -> 67.347C
		require.NotNil(t, frame.TireTemp)
		assert.InDelta(t, 67.347, frame.TireTemp[0], 0.01)
		assert.InDelta(t, 66.704, frame.TireTemp[1], 0.01)
		assert.InDelta(t, 67.551, frame.TireTemp[2], 0.01)
		assert.InDelta(t, 67.551, frame.TireTemp[3], 0.01)

		assert.InDelta(t, -4.641847, *frame.PositionX, 0.001)
		assert.InDelta(t, 369.80182, *frame.PositionY, 0.001)
		assert.InDelta(t, 4243.9, *frame.PositionZ, 0.1)

		assert.InDelta(t, 1.8943167, *frame.Yaw, 0.0001)
		assert.InDelta(t, 0.0639811, *frame.Pitch, 0.0001)
		assert.InDelta(t, -0.007876748, *frame.Roll, 0.0001)

		assert.Equal(t, float32(1.0), *frame.Fuel)
		assert.Equal(t, int32(2), *frame.CarClass)
		assert.Equal(t, int32(3702), *frame.CarIndex)
		assert.Equal(t, int16(0), *frame.LapNumber)

		// Fields Forza doesn't provide should still be set (from packet zero values)
		assert.Nil(t, frame.OilTemp)
		assert.Nil(t, frame.WaterTemp)
		assert.Nil(t, frame.FuelCapacity)
		assert.Nil(t, frame.NumGears)
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
