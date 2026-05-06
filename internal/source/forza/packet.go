package forza

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Packet sizes for each Forza telemetry protocol version.
const (
	PacketSizeV1 = 311 // Forza Motorsport 7
	PacketSizeV2 = 324 // Forza Horizon 4, Forza Horizon 5
	PacketSizeV3 = 331 // Forza Motorsport (2023)
)

// Packet represents decoded Forza UDP telemetry data.
// Fields are the same across all versions; the wire format differs.
type Packet struct {
	IsRaceOn      int32
	TimestampMS   uint32
	EngineMaxRPM  float32
	EngineIdleRPM float32
	CurrentRPM    float32

	AccelerationX float32
	AccelerationY float32
	AccelerationZ float32

	VelocityX float32
	VelocityY float32
	VelocityZ float32

	AngularVelocityX float32
	AngularVelocityY float32
	AngularVelocityZ float32

	Yaw   float32
	Pitch float32
	Roll  float32

	NormalizedSuspensionTravel [4]float32
	TireSlipRatio              [4]float32
	WheelRotationSpeed         [4]float32
	WheelOnRumbleStrip         [4]int32
	WheelInPuddleDepth         [4]float32
	SurfaceRumble              [4]float32
	TireSlipAngle              [4]float32
	TireCombinedSlip           [4]float32
	SuspensionTravelMeters     [4]float32

	CarOrdinal          int32
	CarClass            int32
	CarPerformanceIndex int32
	DrivetrainType      int32
	NumCylinders        int32

	PositionX float32
	PositionY float32
	PositionZ float32

	Speed  float32
	Power  float32
	Torque float32

	TireTemp [4]float32

	Boost            float32
	Fuel             float32
	DistanceTraveled float32
	BestLap          float32
	LastLap          float32
	CurrentLap       float32
	CurrentRaceTime  float32

	LapNumber    uint16
	RacePosition uint8
	Accel        uint8
	Brake        uint8
	Clutch       uint8
	HandBrake    uint8
	Gear         uint8
	Steer        int8

	NormalizedDrivingLine int8
	NormalizedAIBrakeDiff int8
}

func decodePacket(data []byte) (*Packet, error) {
	var posOffset int

	switch len(data) {
	case PacketSizeV1:
		posOffset = 232
	case PacketSizeV2:
		posOffset = 244 // 12-byte gap after NumCylinders
	case PacketSizeV3:
		posOffset = 232 // same layout as V1, extra bytes at end
	default:
		return nil, fmt.Errorf("invalid packet size: %d, expected %d, %d, or %d",
			len(data), PacketSizeV1, PacketSizeV2, PacketSizeV3)
	}

	pkt := &Packet{
		IsRaceOn:      readInt32(data, 0),
		TimestampMS:   binary.LittleEndian.Uint32(data[4:8]),
		EngineMaxRPM:  readFloat32(data, 8),
		EngineIdleRPM: readFloat32(data, 12),
		CurrentRPM:    readFloat32(data, 16),

		AccelerationX: readFloat32(data, 20),
		AccelerationY: readFloat32(data, 24),
		AccelerationZ: readFloat32(data, 28),

		VelocityX: readFloat32(data, 32),
		VelocityY: readFloat32(data, 36),
		VelocityZ: readFloat32(data, 40),

		AngularVelocityX: readFloat32(data, 44),
		AngularVelocityY: readFloat32(data, 48),
		AngularVelocityZ: readFloat32(data, 52),

		Yaw:   readFloat32(data, 56),
		Pitch: readFloat32(data, 60),
		Roll:  readFloat32(data, 64),
	}

	// Per-wheel data [FL, FR, RL, RR] - 4 floats each
	readFloat32x4(data, 68, &pkt.NormalizedSuspensionTravel)
	readFloat32x4(data, 84, &pkt.TireSlipRatio)
	readFloat32x4(data, 100, &pkt.WheelRotationSpeed)
	readInt32x4(data, 116, &pkt.WheelOnRumbleStrip)
	readFloat32x4(data, 132, &pkt.WheelInPuddleDepth)
	readFloat32x4(data, 148, &pkt.SurfaceRumble)
	readFloat32x4(data, 164, &pkt.TireSlipAngle)
	readFloat32x4(data, 180, &pkt.TireCombinedSlip)
	readFloat32x4(data, 196, &pkt.SuspensionTravelMeters)

	// Car identity at offset 212
	pkt.CarOrdinal = readInt32(data, 212)
	pkt.CarClass = readInt32(data, 216)
	pkt.CarPerformanceIndex = readInt32(data, 220)
	pkt.DrivetrainType = readInt32(data, 224)
	pkt.NumCylinders = readInt32(data, 228)

	// Position and remaining fields at version-dependent offset
	off := posOffset

	pkt.PositionX = readFloat32(data, off)
	pkt.PositionY = readFloat32(data, off+4)
	pkt.PositionZ = readFloat32(data, off+8)
	off += 12

	pkt.Speed = readFloat32(data, off)
	pkt.Power = readFloat32(data, off+4)
	pkt.Torque = readFloat32(data, off+8)
	off += 12

	readFloat32x4(data, off, &pkt.TireTemp)
	off += 16

	pkt.Boost = readFloat32(data, off)
	pkt.Fuel = readFloat32(data, off+4)
	pkt.DistanceTraveled = readFloat32(data, off+8)
	pkt.BestLap = readFloat32(data, off+12)
	pkt.LastLap = readFloat32(data, off+16)
	pkt.CurrentLap = readFloat32(data, off+20)
	pkt.CurrentRaceTime = readFloat32(data, off+24)
	off += 28

	pkt.LapNumber = binary.LittleEndian.Uint16(data[off : off+2])
	pkt.RacePosition = data[off+2]
	pkt.Accel = data[off+3]
	pkt.Brake = data[off+4]
	pkt.Clutch = data[off+5]
	pkt.HandBrake = data[off+6]
	pkt.Gear = data[off+7]
	pkt.Steer = int8(data[off+8])
	pkt.NormalizedDrivingLine = int8(data[off+9])
	pkt.NormalizedAIBrakeDiff = int8(data[off+10])

	return pkt, nil
}

func readFloat32(data []byte, offset int) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
}

func readInt32(data []byte, offset int) int32 {
	return int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
}

func readFloat32x4(data []byte, offset int, out *[4]float32) {
	for i := range 4 {
		out[i] = readFloat32(data, offset+i*4)
	}
}

func readInt32x4(data []byte, offset int, out *[4]int32) {
	for i := range 4 {
		out[i] = readInt32(data, offset+i*4)
	}
}
