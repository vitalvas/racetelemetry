package pcars2udp

import (
	"github.com/vitalvas/racetelemetry/internal/model"
)

func val[T any](p *T) T {
	if p != nil {
		return *p
	}

	var zero T

	return zero
}

func valArr4[T any](p *[4]T) [4]T {
	if p != nil {
		return *p
	}

	var zero [4]T

	return zero
}

// convertTelemetry converts a unified frame to pCars2 sTelemetryData UDP packet.
func convertTelemetry(seq uint32, frame *model.TelemetryFrame) []byte {
	buf := make([]byte, telemetryPacketSize)

	encodeHeader(buf, PacketHeader{
		PacketNumber:         seq,
		CategoryPacketNumber: seq,
		PartialPacketIndex:   1,
		PartialPacketNumber:  1,
		PacketType:           PacketTypeCarPhysics,
		PacketVersion:        PacketVersion,
	})

	off := headerSize

	// Viewed participant index (int8)
	buf[off] = 0
	off++

	// Unfiltered throttle (uint8, 0-255)
	buf[off] = uint8(val(frame.Throttle) * 255)
	off++
	// Unfiltered brake (uint8, 0-255)
	buf[off] = uint8(val(frame.Brake) * 255)
	off++
	// Unfiltered steering (int8, -127 to 127)
	buf[off] = uint8(int8(val(frame.Steer) * 127))
	off++
	// Unfiltered clutch (uint8, 0-255)
	buf[off] = uint8(val(frame.Clutch) * 255)
	off++

	// Race state flags (uint8)
	var raceStateFlags uint8
	if val(frame.IsRaceOn) {
		raceStateFlags = 2
	}

	buf[off] = raceStateFlags
	off++

	// Laps in event (uint8)
	buf[off] = uint8(val(frame.TotalLaps))
	off++

	// Padding
	off++

	// Speed (float32, m/s)
	putFloat32(buf, off, val(frame.Speed))
	off += 4

	putFloat32(buf, off, val(frame.EngineRPM))
	off += 4

	putFloat32(buf, off, val(frame.EngineMaxRPM))
	off += 4

	putFloat32(buf, off, val(frame.Throttle))
	off += 4

	putFloat32(buf, off, val(frame.Brake))
	off += 4

	putFloat32(buf, off, val(frame.Clutch))
	off += 4

	putFloat32(buf, off, val(frame.Steer))
	off += 4

	putFloat32(buf, off, val(frame.Fuel))
	off += 4

	putFloat32(buf, off, val(frame.FuelCapacity))
	off += 4

	tireTemp := valArr4(frame.TireTemp)

	// Brake temperature [4]
	for i := 0; i < 4; i++ {
		putFloat32(buf, off, tireTemp[i])
		off += 4
	}

	// Tire temp [4]
	for i := 0; i < 4; i++ {
		putFloat32(buf, off, tireTemp[i])
		off += 4
	}

	// Tire pressure [4]
	for i := 0; i < 4; i++ {
		putFloat32(buf, off, 200.0)
		off += 4
	}

	suspTravel := valArr4(frame.SuspensionTravel)

	for i := 0; i < 4; i++ {
		putFloat32(buf, off, suspTravel[i])
		off += 4
	}

	wheelSpeed := valArr4(frame.WheelSpeed)

	for i := 0; i < 4; i++ {
		putFloat32(buf, off, wheelSpeed[i])
		off += 4
	}

	slipRatio := valArr4(frame.SlipRatio)

	for i := 0; i < 4; i++ {
		putFloat32(buf, off, slipRatio[i])
		off += 4
	}

	// Orientation
	putFloat32(buf, off, val(frame.Yaw))
	off += 4
	putFloat32(buf, off, val(frame.Pitch))
	off += 4
	putFloat32(buf, off, val(frame.Roll))
	off += 4

	// Local velocity
	putFloat32(buf, off, val(frame.VelocityX))
	off += 4
	putFloat32(buf, off, val(frame.VelocityY))
	off += 4
	putFloat32(buf, off, val(frame.VelocityZ))
	off += 4

	// World velocity (approximate from local)
	putFloat32(buf, off, val(frame.VelocityX))
	off += 4
	putFloat32(buf, off, val(frame.VelocityY))
	off += 4
	putFloat32(buf, off, val(frame.VelocityZ))
	off += 4

	// Angular velocity
	putFloat32(buf, off, val(frame.AngularVelocityX))
	off += 4
	putFloat32(buf, off, val(frame.AngularVelocityY))
	off += 4
	putFloat32(buf, off, val(frame.AngularVelocityZ))
	off += 4

	// Local acceleration
	putFloat32(buf, off, val(frame.AccelerationX))
	off += 4
	putFloat32(buf, off, val(frame.AccelerationY))
	off += 4
	putFloat32(buf, off, val(frame.AccelerationZ))
	off += 4

	// World acceleration (approximate from local)
	putFloat32(buf, off, val(frame.AccelerationX))
	off += 4
	putFloat32(buf, off, val(frame.AccelerationY))
	off += 4
	putFloat32(buf, off, val(frame.AccelerationZ))
	off += 4

	// Extents centre
	putFloat32(buf, off, val(frame.PositionX))
	off += 4
	putFloat32(buf, off, val(frame.PositionY))
	off += 4
	putFloat32(buf, off, val(frame.PositionZ))
	off += 4

	// Gear (int8)
	buf[off] = uint8(val(frame.Gear))
	off++

	// Num gears (int8)
	buf[off] = uint8(val(frame.NumGears))
	off++

	// Odometer - not tracked
	off += 4

	putFloat32(buf, off, val(frame.Boost))
	off += 4

	putFloat32(buf, off, val(frame.OilTemp))
	off += 4

	putFloat32(buf, off, val(frame.WaterTemp))
	off += 4

	// Oil pressure (float32, kPa)
	putFloat32(buf, off, val(frame.OilPressure))
	off += 4

	// Water pressure (float32, kPa)
	off += 4

	// Fuel pressure (float32, kPa)
	off += 4

	// Car flags (uint8)
	var carFlags uint8
	if val(frame.IsRaceOn) {
		carFlags |= 0x02 // ENGINE_ACTIVE
	}

	if val(frame.HandBrake) > 0 {
		carFlags |= 0x20 // HANDBRAKE
	}

	buf[off] = carFlags
	off++

	// Engine torque (float32, Nm)
	putFloat32(buf, off, val(frame.Torque))
	off += 4

	// Engine speed (float32)
	putFloat32(buf, off, val(frame.EngineRPM))
	off += 4

	// Wings (float32 x2)
	off += 8

	// HandBrake (float32)
	putFloat32(buf, off, val(frame.HandBrake))

	return buf
}

// convertGameState creates a pCars2 sGameStateData packet.
func convertGameState(seq uint32, frame *model.TelemetryFrame) []byte {
	buf := make([]byte, gameStatePacketSize)

	encodeHeader(buf, PacketHeader{
		PacketNumber:         seq,
		CategoryPacketNumber: seq,
		PartialPacketIndex:   1,
		PartialPacketNumber:  1,
		PacketType:           PacketTypeGameState,
		PacketVersion:        PacketVersion,
	})

	off := headerSize

	if val(frame.IsRaceOn) {
		putUint16(buf, off, 2)
	}

	off += 2

	if val(frame.IsRaceOn) {
		putUint16(buf, off, 5)
	}

	return buf
}

// convertTimings creates a pCars2 sTimingsData packet.
func convertTimings(seq uint32, frame *model.TelemetryFrame) []byte {
	buf := make([]byte, timingsPacketSize)

	encodeHeader(buf, PacketHeader{
		PacketNumber:         seq,
		CategoryPacketNumber: seq,
		PartialPacketIndex:   1,
		PartialPacketNumber:  1,
		PacketType:           PacketTypeTimings,
		PacketVersion:        PacketVersion,
	})

	off := headerSize

	buf[off] = 1
	off++

	off += 4

	putFloat32(buf, off, -1.0)
	off += 4

	off += 8

	putFloat32(buf, off, val(frame.PositionX))
	off += 4
	putFloat32(buf, off, val(frame.PositionY))
	off += 4
	putFloat32(buf, off, val(frame.PositionZ))
	off += 4

	putFloat32(buf, off, val(frame.LapDistance))
	off += 4

	buf[off] = val(frame.RacePosition)
	off++

	lapNum := val(frame.LapNumber)
	if lapNum > 0 {
		buf[off] = uint8(lapNum - 1)
	}

	off++

	buf[off] = uint8(lapNum)
	off++

	off++

	putFloat32(buf, off, val(frame.BestLapTime))
	off += 4

	putFloat32(buf, off, val(frame.LastLapTime))
	off += 4

	putFloat32(buf, off, val(frame.CurrentLapTime))
	off += 4

	if val(frame.IsRaceOn) {
		buf[off] = 2
	}

	return buf
}
