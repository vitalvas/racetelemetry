package pcars1udp

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

// convertTelemetry converts a unified frame to pCars1 sTelemetryData UDP packet.
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

	buf[off] = 0
	off++

	buf[off] = uint8(val(frame.Throttle) * 255)
	off++

	buf[off] = uint8(val(frame.Brake) * 255)
	off++

	buf[off] = uint8(int8(val(frame.Steer) * 127))
	off++

	buf[off] = uint8(val(frame.Clutch) * 255)
	off++

	var raceStateFlags uint8
	if val(frame.IsRaceOn) {
		raceStateFlags = 2
	}

	buf[off] = raceStateFlags
	off++

	buf[off] = uint8(val(frame.TotalLaps))
	off++

	off++

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

	for i := 0; i < 4; i++ {
		putFloat32(buf, off, tireTemp[i])
		off += 4
	}

	for i := 0; i < 4; i++ {
		putFloat32(buf, off, tireTemp[i])
		off += 4
	}

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

	putFloat32(buf, off, val(frame.Yaw))
	off += 4
	putFloat32(buf, off, val(frame.Pitch))
	off += 4
	putFloat32(buf, off, val(frame.Roll))
	off += 4

	putFloat32(buf, off, val(frame.VelocityX))
	off += 4
	putFloat32(buf, off, val(frame.VelocityY))
	off += 4
	putFloat32(buf, off, val(frame.VelocityZ))
	off += 4

	putFloat32(buf, off, val(frame.VelocityX))
	off += 4
	putFloat32(buf, off, val(frame.VelocityY))
	off += 4
	putFloat32(buf, off, val(frame.VelocityZ))
	off += 4

	putFloat32(buf, off, val(frame.AngularVelocityX))
	off += 4
	putFloat32(buf, off, val(frame.AngularVelocityY))
	off += 4
	putFloat32(buf, off, val(frame.AngularVelocityZ))
	off += 4

	putFloat32(buf, off, val(frame.AccelerationX))
	off += 4
	putFloat32(buf, off, val(frame.AccelerationY))
	off += 4
	putFloat32(buf, off, val(frame.AccelerationZ))
	off += 4

	putFloat32(buf, off, val(frame.AccelerationX))
	off += 4
	putFloat32(buf, off, val(frame.AccelerationY))
	off += 4
	putFloat32(buf, off, val(frame.AccelerationZ))
	off += 4

	putFloat32(buf, off, val(frame.PositionX))
	off += 4
	putFloat32(buf, off, val(frame.PositionY))
	off += 4
	putFloat32(buf, off, val(frame.PositionZ))
	off += 4

	buf[off] = uint8(val(frame.Gear))
	off++

	buf[off] = uint8(val(frame.NumGears))
	off++

	off += 4

	putFloat32(buf, off, val(frame.Boost))
	off += 4

	putFloat32(buf, off, val(frame.OilTemp))
	off += 4

	putFloat32(buf, off, val(frame.WaterTemp))
	off += 4

	off += 12

	var carFlags uint8
	if val(frame.IsRaceOn) {
		carFlags |= 0x02
	}

	buf[off] = carFlags

	return buf
}

// convertGameState creates a pCars1 sGameStateData packet.
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

// convertTimings creates a pCars1 sTimingsData packet.
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
