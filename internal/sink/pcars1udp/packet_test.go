package pcars1udp

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeHeader(t *testing.T) {
	buf := make([]byte, headerSize)

	h := PacketHeader{
		PacketNumber:         100,
		CategoryPacketNumber: 50,
		PartialPacketIndex:   1,
		PartialPacketNumber:  1,
		PacketType:           PacketTypeCarPhysics,
		PacketVersion:        PacketVersion,
	}

	encodeHeader(buf, h)

	assert.Equal(t, uint32(100), binary.LittleEndian.Uint32(buf[0:4]))
	assert.Equal(t, uint32(50), binary.LittleEndian.Uint32(buf[4:8]))
	assert.Equal(t, uint8(1), buf[8])
	assert.Equal(t, uint8(1), buf[9])
	assert.Equal(t, uint8(PacketTypeCarPhysics), buf[10])
	assert.Equal(t, uint8(PacketVersion), buf[11])
}

func TestPutFloat32(t *testing.T) {
	buf := make([]byte, 8)
	putFloat32(buf, 0, 3.14)

	bits := binary.LittleEndian.Uint32(buf[0:4])
	assert.Equal(t, math.Float32bits(3.14), bits)
}
