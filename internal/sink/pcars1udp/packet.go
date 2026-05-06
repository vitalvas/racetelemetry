package pcars1udp

import (
	"encoding/binary"
	"math"
)

// Packet types as defined by pCars1 UDP protocol.
const (
	PacketTypeCarPhysics     = 0
	PacketTypeRaceDefinition = 1
	PacketTypeParticipants   = 2
	PacketTypeTimings        = 3
	PacketTypeGameState      = 4

	PacketVersion = 1

	headerSize          = 12
	telemetryPacketSize = 538
	gameStatePacketSize = 16
	timingsPacketSize   = 993
)

// PacketHeader is the 12-byte common header for all pCars1 UDP packets.
type PacketHeader struct {
	PacketNumber         uint32
	CategoryPacketNumber uint32
	PartialPacketIndex   uint8
	PartialPacketNumber  uint8
	PacketType           uint8
	PacketVersion        uint8
}

func encodeHeader(buf []byte, h PacketHeader) {
	binary.LittleEndian.PutUint32(buf[0:4], h.PacketNumber)
	binary.LittleEndian.PutUint32(buf[4:8], h.CategoryPacketNumber)
	buf[8] = h.PartialPacketIndex
	buf[9] = h.PartialPacketNumber
	buf[10] = h.PacketType
	buf[11] = h.PacketVersion
}

func putFloat32(buf []byte, offset int, v float32) {
	binary.LittleEndian.PutUint32(buf[offset:offset+4], math.Float32bits(v))
}

func putUint16(buf []byte, offset int, v uint16) {
	binary.LittleEndian.PutUint16(buf[offset:offset+2], v)
}
