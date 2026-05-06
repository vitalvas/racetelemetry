package pcars2udp

import (
	"github.com/vitalvas/racetelemetry/internal/model"
)

// convertTelemetry converts a unified frame to pCars2 sTelemetryData UDP packet.
//
// pCars2 sTelemetryData layout (after 12-byte header):
// The packet contains car physics data that CrewChief uses for voice callouts.
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

	// Viewed participant index (int8) - always 0 (player)
	buf[off] = 0
	off++

	// Unfiltered throttle (uint8, 0-255)
	buf[off] = uint8(frame.Throttle * 255)
	off++
	// Unfiltered brake (uint8, 0-255)
	buf[off] = uint8(frame.Brake * 255)
	off++
	// Unfiltered steering (int8, -127 to 127)
	buf[off] = uint8(int8(frame.Steer * 127))
	off++
	// Unfiltered clutch (uint8, 0-255)
	buf[off] = uint8(frame.Clutch * 255)
	off++

	// Race state flags (uint8)
	var raceStateFlags uint8
	if frame.IsRaceOn {
		raceStateFlags = 2 // RACING
	}
	buf[off] = raceStateFlags
	off++

	// Laps in event (uint8)
	buf[off] = uint8(frame.TotalLaps)
	off++

	// Padding/reserved (1 byte)
	off++

	// Speed (float32, m/s)
	putFloat32(buf, off, frame.Speed)
	off += 4

	// RPM (float32)
	putFloat32(buf, off, frame.EngineRPM)
	off += 4

	// Max RPM (float32)
	putFloat32(buf, off, frame.EngineMaxRPM)
	off += 4

	// Throttle (float32, 0.0-1.0)
	putFloat32(buf, off, frame.Throttle)
	off += 4

	// Brake (float32, 0.0-1.0)
	putFloat32(buf, off, frame.Brake)
	off += 4

	// Clutch (float32, 0.0-1.0)
	putFloat32(buf, off, frame.Clutch)
	off += 4

	// Steering (float32, -1.0 to 1.0)
	putFloat32(buf, off, frame.Steer)
	off += 4

	// Fuel level (float32, 0.0-1.0)
	putFloat32(buf, off, frame.Fuel)
	off += 4

	// Fuel capacity (float32)
	putFloat32(buf, off, frame.FuelCapacity)
	off += 4

	// Brake temperature [4] (float32, Celsius)
	// We approximate from tire temp since source doesn't provide brake temp separately
	for i := 0; i < 4; i++ {
		putFloat32(buf, off, frame.TireTemp[i])
		off += 4
	}

	// Tire temp [4] (float32, Celsius)
	for i := 0; i < 4; i++ {
		putFloat32(buf, off, frame.TireTemp[i])
		off += 4
	}

	// Tire pressure (set to typical 200 kPa)
	for i := 0; i < 4; i++ {
		putFloat32(buf, off, 200.0)
		off += 4
	}

	// Suspension travel [4] (float32, meters)
	for i := 0; i < 4; i++ {
		putFloat32(buf, off, frame.SuspensionTravel[i])
		off += 4
	}

	// Wheel speed [4] (float32, rad/s)
	for i := 0; i < 4; i++ {
		putFloat32(buf, off, frame.WheelSpeed[i])
		off += 4
	}

	// Tire slip speed [4] (float32)
	for i := 0; i < 4; i++ {
		putFloat32(buf, off, frame.SlipRatio[i])
		off += 4
	}

	// Orientation (float32 x3)
	putFloat32(buf, off, frame.Yaw)
	off += 4
	putFloat32(buf, off, frame.Pitch)
	off += 4
	putFloat32(buf, off, frame.Roll)
	off += 4

	// Local velocity (float32 x3)
	putFloat32(buf, off, frame.VelocityX)
	off += 4
	putFloat32(buf, off, frame.VelocityY)
	off += 4
	putFloat32(buf, off, frame.VelocityZ)
	off += 4

	// World velocity (float32 x3) - approximate from local
	putFloat32(buf, off, frame.VelocityX)
	off += 4
	putFloat32(buf, off, frame.VelocityY)
	off += 4
	putFloat32(buf, off, frame.VelocityZ)
	off += 4

	// Angular velocity (float32 x3)
	putFloat32(buf, off, frame.AngularVelocityX)
	off += 4
	putFloat32(buf, off, frame.AngularVelocityY)
	off += 4
	putFloat32(buf, off, frame.AngularVelocityZ)
	off += 4

	// Local acceleration (float32 x3)
	putFloat32(buf, off, frame.AccelerationX)
	off += 4
	putFloat32(buf, off, frame.AccelerationY)
	off += 4
	putFloat32(buf, off, frame.AccelerationZ)
	off += 4

	// World acceleration (float32 x3) - approximate from local
	putFloat32(buf, off, frame.AccelerationX)
	off += 4
	putFloat32(buf, off, frame.AccelerationY)
	off += 4
	putFloat32(buf, off, frame.AccelerationZ)
	off += 4

	// Extents centre (float32 x3)
	putFloat32(buf, off, frame.PositionX)
	off += 4
	putFloat32(buf, off, frame.PositionY)
	off += 4
	putFloat32(buf, off, frame.PositionZ)
	off += 4

	// Gear (int8) - pCars2 uses -1=R, 0=N, 1+=forward (matches our unified model)
	buf[off] = uint8(frame.Gear)
	off++

	// Num gears (int8)
	buf[off] = uint8(frame.NumGears)
	off++

	// Odometer (float32) - not tracked, sources only provide per-lap distance
	// putFloat32(buf, off, 0) — zero-initialized
	off += 4

	// Boost amount (float32)
	putFloat32(buf, off, frame.Boost)
	off += 4

	// Oil temp (float32, Celsius)
	putFloat32(buf, off, frame.OilTemp)
	off += 4

	// Water temp (float32, Celsius)
	putFloat32(buf, off, frame.WaterTemp)
	off += 4

	// Oil pressure (float32, kPa)
	off += 4

	// Water pressure (float32, kPa)
	off += 4

	// Fuel pressure (float32, kPa)
	off += 4

	// Car flags (uint8)
	var carFlags uint8
	if frame.IsRaceOn {
		carFlags |= 0x02 // ENGINE_ACTIVE
	}
	buf[off] = carFlags
	off++

	// Engine torque (float32, Nm)
	putFloat32(buf, off, frame.Torque)
	off += 4

	// Engine speed (float32)
	putFloat32(buf, off, frame.EngineRPM)
	// remaining bytes are zero-initialized (wings, handbrake, etc)

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

	// Game state (uint16) - INGAME_PLAYING = 2
	if frame.IsRaceOn {
		putUint16(buf, off, 2)
	}
	off += 2

	// Session state (uint16) - RACE = 5
	if frame.IsRaceOn {
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

	// Num participants (int8)
	buf[off] = 1
	off++

	// Participants changed timestamp (uint32)
	off += 4

	// Event time remaining (float32, -1 for no limit)
	putFloat32(buf, off, -1.0)
	off += 4

	// Split time ahead (float32)
	off += 4
	// Split time behind (float32)
	off += 4

	// Participant info for player (index 0)
	// World position [3] (float32)
	putFloat32(buf, off, frame.PositionX)
	off += 4
	putFloat32(buf, off, frame.PositionY)
	off += 4
	putFloat32(buf, off, frame.PositionZ)
	off += 4

	// Current lap distance (float32)
	putFloat32(buf, off, frame.LapDistance)
	off += 4

	// Race position (uint8)
	buf[off] = frame.RacePosition
	off++

	// Laps completed (uint8)
	if frame.LapNumber > 0 {
		buf[off] = uint8(frame.LapNumber - 1)
	}

	off++

	// Current lap (uint8)
	buf[off] = uint8(frame.LapNumber)
	off++

	// Current sector (int8) - 0 = sector 1
	off++

	// Fastest lap time (float32)
	putFloat32(buf, off, frame.BestLapTime)
	off += 4

	// Last lap time (float32)
	putFloat32(buf, off, frame.LastLapTime)
	off += 4

	// Current lap time (float32)
	putFloat32(buf, off, frame.CurrentLapTime)
	off += 4

	// Race state per participant (uint8)
	if frame.IsRaceOn {
		buf[off] = 2 // RACING
	}

	return buf
}
