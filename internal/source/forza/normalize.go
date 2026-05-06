package forza

import (
	"time"

	"github.com/vitalvas/racetelemetry/internal/model"
)

func normalize(pkt *Packet) model.TelemetryFrame {
	frame := model.TelemetryFrame{
		Timestamp: time.Now(),
		IsRaceOn:  pkt.IsRaceOn > 0,

		EngineRPM:     pkt.CurrentRPM,
		EngineMaxRPM:  pkt.EngineMaxRPM,
		EngineIdleRPM: pkt.EngineIdleRPM,

		Gear:      normalizeGear(pkt.Gear),
		Speed:     pkt.Speed,
		Throttle:  float32(pkt.Accel) / 255.0,
		Brake:     float32(pkt.Brake) / 255.0,
		Clutch:    float32(pkt.Clutch) / 255.0,
		HandBrake: float32(pkt.HandBrake) / 255.0,
		Steer:     float32(pkt.Steer) / 127.0,

		AccelerationX:    pkt.AccelerationX,
		AccelerationY:    pkt.AccelerationY,
		AccelerationZ:    pkt.AccelerationZ,
		VelocityX:        pkt.VelocityX,
		VelocityY:        pkt.VelocityY,
		VelocityZ:        pkt.VelocityZ,
		AngularVelocityX: pkt.AngularVelocityX,
		AngularVelocityY: pkt.AngularVelocityY,
		AngularVelocityZ: pkt.AngularVelocityZ,

		Yaw:   pkt.Yaw,
		Pitch: pkt.Pitch,
		Roll:  pkt.Roll,

		PositionX: pkt.PositionX,
		PositionY: pkt.PositionY,
		PositionZ: pkt.PositionZ,

		Power:  pkt.Power,
		Torque: pkt.Torque,
		Boost:  pkt.Boost,

		Fuel: pkt.Fuel,

		LapNumber:      int16(pkt.LapNumber),
		RacePosition:   pkt.RacePosition,
		BestLapTime:    pkt.BestLap,
		LastLapTime:    pkt.LastLap,
		CurrentLapTime: pkt.CurrentLap,
		LapDistance:    pkt.DistanceTraveled,

		CarClass: pkt.CarClass,
		CarIndex: pkt.CarOrdinal,
	}

	// Convert tire temperatures from Fahrenheit to Celsius
	for i := 0; i < 4; i++ {
		frame.TireTemp[i] = fahrenheitToCelsius(pkt.TireTemp[i])
	}

	frame.SuspensionTravel = pkt.SuspensionTravelMeters
	frame.WheelSpeed = pkt.WheelRotationSpeed
	frame.SlipRatio = pkt.TireSlipRatio
	frame.SlipAngle = pkt.TireSlipAngle

	return frame
}

// normalizeGear converts Forza gear encoding to unified encoding.
// Forza: 0=Reverse, 1=Neutral, 2=1st, 3=2nd, ...
// Unified: -1=Reverse, 0=Neutral, 1=1st, 2=2nd, ...
func normalizeGear(g uint8) int8 {
	switch g {
	case 0:
		return -1
	case 1:
		return 0
	default:
		return int8(g) - 1
	}
}

func fahrenheitToCelsius(f float32) float32 {
	return (f - 32.0) * 5.0 / 9.0
}
