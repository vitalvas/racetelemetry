package gt7

import (
	"time"

	"github.com/vitalvas/racetelemetry/internal/model"
)

func normalize(pkt *Packet) model.TelemetryFrame {
	isRaceOn := pkt.StatusFlags&0x01 != 0
	gear := normalizeGear(pkt.CurrentGear)
	maxRPM := float32(pkt.RPMRevLimiter)
	throttle := float32(pkt.Throttle) / 255.0
	brake := float32(pkt.Brake) / 255.0
	boost := pkt.Boost - 1.0
	bestLap := float32(pkt.BestLapTime) / 1000.0
	lastLap := float32(pkt.LastLapTime) / 1000.0
	racePos := uint8(pkt.CurrentPosition)

	frame := model.TelemetryFrame{
		Timestamp: time.Now(),
		IsRaceOn:  &isRaceOn,

		EngineRPM:    &pkt.RPM,
		EngineMaxRPM: &maxRPM,

		Gear:     &gear,
		Speed:    &pkt.CarSpeed,
		Throttle: &throttle,
		Brake:    &brake,
		Clutch:   &pkt.Clutch,
		Steer:    &pkt.SteeringAngle,

		VelocityX:        &pkt.VelocityX,
		VelocityY:        &pkt.VelocityY,
		VelocityZ:        &pkt.VelocityZ,
		AngularVelocityX: &pkt.AngularVelocityX,
		AngularVelocityY: &pkt.AngularVelocityY,
		AngularVelocityZ: &pkt.AngularVelocityZ,

		Yaw:   &pkt.RotationYaw,
		Pitch: &pkt.RotationPitch,
		Roll:  &pkt.RotationRoll,

		PositionX: &pkt.PositionX,
		PositionY: &pkt.PositionY,
		PositionZ: &pkt.PositionZ,

		Boost: &boost,

		FuelCapacity: &pkt.FuelCapacity,

		OilTemp:   &pkt.OilTemp,
		WaterTemp: &pkt.WaterTemp,

		TireTemp:         &pkt.TireTemp,
		SuspensionTravel: &pkt.Suspension,
		WheelSpeed:       &pkt.TyreAngularSpeed,

		LapNumber:    &pkt.CurrentLap,
		TotalLaps:    &pkt.TotalLaps,
		RacePosition: &racePos,
		BestLapTime:  &bestLap,
		LastLapTime:  &lastLap,

		CarIndex:   &pkt.CarID,
		GearRatios: &pkt.GearRatios,
	}

	if pkt.FuelCapacity > 0 {
		fuel := pkt.CurrentFuel / pkt.FuelCapacity
		frame.Fuel = &fuel
	}

	if pkt.CurrentLapTimeMS != 0 {
		lapTime := float32(pkt.CurrentLapTimeMS) / 1000.0
		frame.CurrentLapTime = &lapTime
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
