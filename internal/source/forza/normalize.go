package forza

import (
	"time"

	"github.com/vitalvas/racetelemetry/internal/model"
)

func normalize(pkt *Packet) model.TelemetryFrame {
	isRaceOn := pkt.IsRaceOn > 0
	gear := normalizeGear(pkt.Gear)
	throttle := float32(pkt.Accel) / 255.0
	brake := float32(pkt.Brake) / 255.0
	clutch := float32(pkt.Clutch) / 255.0
	handBrake := float32(pkt.HandBrake) / 255.0
	steer := float32(pkt.Steer) / 127.0

	tireTemp := [4]float32{
		fahrenheitToCelsius(pkt.TireTemp[0]),
		fahrenheitToCelsius(pkt.TireTemp[1]),
		fahrenheitToCelsius(pkt.TireTemp[2]),
		fahrenheitToCelsius(pkt.TireTemp[3]),
	}

	wheelOnRumble := [4]bool{
		pkt.WheelOnRumbleStrip[0] > 0,
		pkt.WheelOnRumbleStrip[1] > 0,
		pkt.WheelOnRumbleStrip[2] > 0,
		pkt.WheelOnRumbleStrip[3] > 0,
	}

	lapNumber := int16(pkt.LapNumber)

	return model.TelemetryFrame{
		Timestamp:       time.Now(),
		IsRaceOn:        &isRaceOn,
		CurrentRaceTime: &pkt.CurrentRaceTime,

		EngineRPM:     &pkt.CurrentRPM,
		EngineMaxRPM:  &pkt.EngineMaxRPM,
		EngineIdleRPM: &pkt.EngineIdleRPM,
		NumCylinders:  &pkt.NumCylinders,

		Gear:           &gear,
		Speed:          &pkt.Speed,
		Throttle:       &throttle,
		Brake:          &brake,
		Clutch:         &clutch,
		HandBrake:      &handBrake,
		Steer:          &steer,
		DrivetrainType: &pkt.DrivetrainType,

		AccelerationX:    &pkt.AccelerationX,
		AccelerationY:    &pkt.AccelerationY,
		AccelerationZ:    &pkt.AccelerationZ,
		VelocityX:        &pkt.VelocityX,
		VelocityY:        &pkt.VelocityY,
		VelocityZ:        &pkt.VelocityZ,
		AngularVelocityX: &pkt.AngularVelocityX,
		AngularVelocityY: &pkt.AngularVelocityY,
		AngularVelocityZ: &pkt.AngularVelocityZ,

		Yaw:   &pkt.Yaw,
		Pitch: &pkt.Pitch,
		Roll:  &pkt.Roll,

		PositionX: &pkt.PositionX,
		PositionY: &pkt.PositionY,
		PositionZ: &pkt.PositionZ,

		Power:  &pkt.Power,
		Torque: &pkt.Torque,
		Boost:  &pkt.Boost,

		Fuel: &pkt.Fuel,

		TireTemp:                   &tireTemp,
		SuspensionTravel:           &pkt.SuspensionTravelMeters,
		NormalizedSuspensionTravel: &pkt.NormalizedSuspensionTravel,
		WheelSpeed:                 &pkt.WheelRotationSpeed,
		WheelOnRumbleStrip:         &wheelOnRumble,
		WheelInPuddleDepth:         &pkt.WheelInPuddleDepth,
		SurfaceRumble:              &pkt.SurfaceRumble,
		SlipRatio:                  &pkt.TireSlipRatio,
		SlipAngle:                  &pkt.TireSlipAngle,
		TireCombinedSlip:           &pkt.TireCombinedSlip,

		LapNumber:      &lapNumber,
		RacePosition:   &pkt.RacePosition,
		BestLapTime:    &pkt.BestLap,
		LastLapTime:    &pkt.LastLap,
		CurrentLapTime: &pkt.CurrentLap,
		LapDistance:    &pkt.DistanceTraveled,

		CarClass:            &pkt.CarClass,
		CarIndex:            &pkt.CarOrdinal,
		CarPerformanceIndex: &pkt.CarPerformanceIndex,

		NormalizedDrivingLine: &pkt.NormalizedDrivingLine,
		NormalizedAIBrakeDiff: &pkt.NormalizedAIBrakeDiff,
	}
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
