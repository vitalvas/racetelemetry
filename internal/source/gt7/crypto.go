package gt7

import (
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/salsa20"
)

const (
	magicNumber = 0x47375330 // "G7S0" - GT7 and GT Sport
	ivOffset    = 0x40
)

// IV XOR seeds for each packet format.
const (
	IVSeedFormatA uint32 = 0xDEADBEAF // Standard
	IVSeedFormatB uint32 = 0xDEADBEEF // Addendum1
	IVSeedFormatT uint32 = 0x55FABB4F // Addendum2 (heartbeat "~")
	IVSeedFormatC uint32 = 0xDEADBEEF // Addendum3
)

const salsa20KeySource = "Simulator Interface Packet GT7 ver 0.0"

func getSalsa20Key() *[32]byte {
	var key [32]byte
	copy(key[:], salsa20KeySource)

	return &key
}

func decrypt(data []byte, ivSeed uint32) ([]byte, error) {
	if len(data) < ivOffset+4 {
		return nil, errors.New("packet too short for IV extraction")
	}

	iv1 := binary.LittleEndian.Uint32(data[ivOffset : ivOffset+4])
	iv2 := iv1 ^ ivSeed

	var nonce [8]byte

	binary.LittleEndian.PutUint32(nonce[0:4], iv2)
	binary.LittleEndian.PutUint32(nonce[4:8], iv1)

	out := make([]byte, len(data))
	salsa20.XORKeyStream(out, data, nonce[:], getSalsa20Key())

	magic := binary.LittleEndian.Uint32(out[0:4])
	if magic != magicNumber {
		return nil, errors.New("invalid magic after decryption")
	}

	return out, nil
}
