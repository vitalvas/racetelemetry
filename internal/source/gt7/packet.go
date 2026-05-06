package gt7

import (
	"encoding/binary"
	"errors"
	"math"
)

// Minimum packet sizes for each format.
const (
	PacketSizeStandard  = 296 // Format "A"
	PacketSizeAddendum1 = 316 // Format "B"
	PacketSizeAddendum2 = 344 // Format "~"
	PacketSizeAddendum3 = 368 // Format "C"
)

// Packet represents the decrypted GT telemetry data.
// Supports GT7, GT Sport, and GT6.
type Packet struct {
	// Standard fields (format "A", 296 bytes)
	PositionX float32
	PositionY float32
	PositionZ float32

	VelocityX float32
	VelocityY float32
	VelocityZ float32

	RotationPitch float32
	RotationYaw   float32
	RotationRoll  float32

	Heading float32 // orientation to North, 0.0=south, 1.0=north

	AngularVelocityX float32
	AngularVelocityY float32
	AngularVelocityZ float32

	RideHeight float32
	RPM        float32

	CurrentFuel  float32
	FuelCapacity float32
	CarSpeed     float32
	Boost        float32
	OilPressure  float32
	WaterTemp    float32
	OilTemp      float32

	TireTemp [4]float32 // [FL, FR, RL, RR] Celsius

	PackageID int32

	CurrentLap      int16
	TotalLaps       int16
	BestLapTime     int32  // ms
	LastLapTime     int32  // ms
	TimeOnTrack     uint32 // ms
	CurrentPosition int16
	TotalPositions  int16

	RPMRevWarning     uint16
	RPMRevLimiter     uint16
	EstimatedTopSpeed int16

	StatusFlags uint16

	CurrentGear   uint8
	SuggestedGear uint8
	Throttle      uint8 // 0-255
	Brake         uint8 // 0-255

	RoadPlaneX    float32
	RoadPlaneY    float32
	RoadPlaneZ    float32
	RoadPlaneDist float32

	TyreAngularSpeed [4]float32 // rad/s
	TyreRadius       [4]float32 // meters
	Suspension       [4]float32

	Clutch         float32
	ClutchEngaged  float32
	RPMAfterClutch float32

	TopSpeedRatio float32
	GearRatios    [8]float32

	CarID int32

	// Addendum1 fields (format "B", 316 bytes)
	SteeringAngle    float32 // radians
	SteeringVelocity float32 // radians per second

	// Addendum3 fields (format "C", 368 bytes)
	CurrentLapTimeMS int32 // milliseconds
}

func decodePacket(data []byte) (*Packet, error) {
	if len(data) < PacketSizeStandard {
		return nil, errors.New("gt7 packet too short")
	}

	p := &Packet{
		PositionX: readFloat32(data, 0x04),
		PositionY: readFloat32(data, 0x08),
		PositionZ: readFloat32(data, 0x0C),

		VelocityX: readFloat32(data, 0x10),
		VelocityY: readFloat32(data, 0x14),
		VelocityZ: readFloat32(data, 0x18),

		RotationPitch: readFloat32(data, 0x1C),
		RotationYaw:   readFloat32(data, 0x20),
		RotationRoll:  readFloat32(data, 0x24),

		Heading: readFloat32(data, 0x28),

		AngularVelocityX: readFloat32(data, 0x2C),
		AngularVelocityY: readFloat32(data, 0x30),
		AngularVelocityZ: readFloat32(data, 0x34),

		RideHeight: readFloat32(data, 0x38),
		RPM:        readFloat32(data, 0x3C),

		CurrentFuel:  readFloat32(data, 0x44),
		FuelCapacity: readFloat32(data, 0x48),
		CarSpeed:     readFloat32(data, 0x4C),
		Boost:        readFloat32(data, 0x50),
		OilPressure:  readFloat32(data, 0x54),
		WaterTemp:    readFloat32(data, 0x58),
		OilTemp:      readFloat32(data, 0x5C),

		TireTemp: [4]float32{
			readFloat32(data, 0x60),
			readFloat32(data, 0x64),
			readFloat32(data, 0x68),
			readFloat32(data, 0x6C),
		},

		PackageID: readInt32(data, 0x70),

		CurrentLap:      readInt16(data, 0x74),
		TotalLaps:       readInt16(data, 0x76),
		BestLapTime:     readInt32(data, 0x78),
		LastLapTime:     readInt32(data, 0x7C),
		TimeOnTrack:     binary.LittleEndian.Uint32(data[0x80:0x84]),
		CurrentPosition: readInt16(data, 0x84),
		TotalPositions:  readInt16(data, 0x86),

		RPMRevWarning:     binary.LittleEndian.Uint16(data[0x88:0x8A]),
		RPMRevLimiter:     binary.LittleEndian.Uint16(data[0x8A:0x8C]),
		EstimatedTopSpeed: readInt16(data, 0x8C),
		StatusFlags:       binary.LittleEndian.Uint16(data[0x8E:0x90]),

		CurrentGear:   data[0x90] & 0x0F,
		SuggestedGear: (data[0x90] >> 4) & 0x0F,
		Throttle:      data[0x91],
		Brake:         data[0x92],

		RoadPlaneX:    readFloat32(data, 0x94),
		RoadPlaneY:    readFloat32(data, 0x98),
		RoadPlaneZ:    readFloat32(data, 0x9C),
		RoadPlaneDist: readFloat32(data, 0xA0),

		TyreAngularSpeed: [4]float32{
			readFloat32(data, 0xA4),
			readFloat32(data, 0xA8),
			readFloat32(data, 0xAC),
			readFloat32(data, 0xB0),
		},

		TyreRadius: [4]float32{
			readFloat32(data, 0xB4),
			readFloat32(data, 0xB8),
			readFloat32(data, 0xBC),
			readFloat32(data, 0xC0),
		},

		Suspension: [4]float32{
			readFloat32(data, 0xC4),
			readFloat32(data, 0xC8),
			readFloat32(data, 0xCC),
			readFloat32(data, 0xD0),
		},

		Clutch:         readFloat32(data, 0xF4),
		ClutchEngaged:  readFloat32(data, 0xF8),
		RPMAfterClutch: readFloat32(data, 0xFC),

		TopSpeedRatio: readFloat32(data, 0x100),

		GearRatios: [8]float32{
			readFloat32(data, 0x104),
			readFloat32(data, 0x108),
			readFloat32(data, 0x10C),
			readFloat32(data, 0x110),
			readFloat32(data, 0x114),
			readFloat32(data, 0x118),
			readFloat32(data, 0x11C),
			readFloat32(data, 0x120),
		},

		CarID: readInt32(data, 0x124),
	}

	// Addendum1 fields (format "B")
	if len(data) >= PacketSizeAddendum1 {
		p.SteeringAngle = readFloat32(data, 0x128)
		p.SteeringVelocity = readFloat32(data, 0x12C)
	}

	// Addendum3 fields (format "C")
	if len(data) >= PacketSizeAddendum3 {
		p.CurrentLapTimeMS = readInt32(data, 0x15C)
	}

	return p, nil
}

func readFloat32(data []byte, offset int) float32 {
	bits := binary.LittleEndian.Uint32(data[offset : offset+4])
	return math.Float32frombits(bits)
}

func readInt32(data []byte, offset int) int32 {
	return int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
}

func readInt16(data []byte, offset int) int16 {
	return int16(binary.LittleEndian.Uint16(data[offset : offset+2]))
}
