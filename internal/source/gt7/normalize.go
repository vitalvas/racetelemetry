package gt7

import (
	"time"

	"github.com/vitalvas/racetelemetry/internal/model"
)

func normalize(pkt *Packet) model.TelemetryFrame {
	frame := model.TelemetryFrame{
		Timestamp: time.Now(),
		IsRaceOn:  pkt.StatusFlags&0x01 != 0,

		EngineRPM:    pkt.RPM,
		EngineMaxRPM: float32(pkt.RPMRevLimiter),

		Gear:     normalizeGear(pkt.CurrentGear),
		Speed:    pkt.CarSpeed,
		Throttle: float32(pkt.Throttle) / 255.0,
		Brake:    float32(pkt.Brake) / 255.0,
		Clutch:   pkt.Clutch,
		Steer:    pkt.SteeringAngle,

		VelocityX:        pkt.VelocityX,
		VelocityY:        pkt.VelocityY,
		VelocityZ:        pkt.VelocityZ,
		AngularVelocityX: pkt.AngularVelocityX,
		AngularVelocityY: pkt.AngularVelocityY,
		AngularVelocityZ: pkt.AngularVelocityZ,

		Yaw:   pkt.RotationYaw,
		Pitch: pkt.RotationPitch,
		Roll:  pkt.RotationRoll,

		PositionX: pkt.PositionX,
		PositionY: pkt.PositionY,
		PositionZ: pkt.PositionZ,

		Boost: pkt.Boost - 1.0,

		FuelCapacity: pkt.FuelCapacity,

		OilTemp:   pkt.OilTemp,
		WaterTemp: pkt.WaterTemp,

		TireTemp:         pkt.TireTemp,
		SuspensionTravel: pkt.Suspension,
		WheelSpeed:       pkt.TyreAngularSpeed,

		LapNumber:      pkt.CurrentLap,
		TotalLaps:      pkt.TotalLaps,
		RacePosition:   uint8(pkt.CurrentPosition),
		BestLapTime:    float32(pkt.BestLapTime) / 1000.0,
		LastLapTime:    float32(pkt.LastLapTime) / 1000.0,
		CurrentLapTime: float32(pkt.CurrentLapTimeMS) / 1000.0,

		CarIndex:   pkt.CarID,
		GearRatios: pkt.GearRatios,
	}

	if pkt.FuelCapacity > 0 {
		frame.Fuel = pkt.CurrentFuel / pkt.FuelCapacity
	}

	return frame
}

// normalizeGear converts GT gear encoding to unified encoding.
// 0=Reverse, 1=1st, 2=2nd, ..., 15(0xF)=Neutral
func normalizeGear(g uint8) int8 {
	if g == 0 {
		return -1
	}

	if g == 15 {
		return 0
	}

	return int8(g)
}
