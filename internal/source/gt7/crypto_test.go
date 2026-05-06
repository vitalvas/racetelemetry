package gt7

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/salsa20"
)

func TestDecrypt(t *testing.T) {
	seeds := []struct {
		name   string
		ivSeed uint32
	}{
		{"format A", IVSeedFormatA},
		{"format B", IVSeedFormatB},
		{"format C", IVSeedFormatC},
		{"format tilde", IVSeedFormatT},
	}

	for _, seed := range seeds {
		t.Run(fmt.Sprintf("valid packet %s", seed.name), func(t *testing.T) {
			ivValue := uint32(42)
			iv2 := ivValue ^ seed.ivSeed

			var nonce [8]byte
			binary.LittleEndian.PutUint32(nonce[0:4], iv2)
			binary.LittleEndian.PutUint32(nonce[4:8], ivValue)

			cleartext := make([]byte, 0x170)
			binary.LittleEndian.PutUint32(cleartext[0:4], magicNumber)

			encrypted := make([]byte, len(cleartext))
			salsa20.XORKeyStream(encrypted, cleartext, nonce[:], getSalsa20Key())
			binary.LittleEndian.PutUint32(encrypted[ivOffset:ivOffset+4], ivValue)

			result, err := decrypt(encrypted, seed.ivSeed)
			require.NoError(t, err)

			magic := binary.LittleEndian.Uint32(result[0:4])
			assert.Equal(t, uint32(magicNumber), magic)
		})
	}

	t.Run("packet too short", func(t *testing.T) {
		_, err := decrypt(make([]byte, 10), IVSeedFormatC)
		assert.Error(t, err)
	})

	t.Run("invalid magic", func(t *testing.T) {
		data := make([]byte, 0x170)
		for i := range data {
			data[i] = byte(i)
		}

		_, err := decrypt(data, IVSeedFormatC)
		assert.Error(t, err)
	})

	t.Run("wrong iv seed fails", func(t *testing.T) {
		ivValue := uint32(42)
		iv2 := ivValue ^ IVSeedFormatA

		var nonce [8]byte
		binary.LittleEndian.PutUint32(nonce[0:4], iv2)
		binary.LittleEndian.PutUint32(nonce[4:8], ivValue)

		cleartext := make([]byte, 0x170)
		binary.LittleEndian.PutUint32(cleartext[0:4], magicNumber)

		encrypted := make([]byte, len(cleartext))
		salsa20.XORKeyStream(encrypted, cleartext, nonce[:], getSalsa20Key())
		binary.LittleEndian.PutUint32(encrypted[ivOffset:ivOffset+4], ivValue)

		// Decrypt with wrong seed should fail magic check
		_, err := decrypt(encrypted, IVSeedFormatC)
		assert.Error(t, err)
	})
}

func BenchmarkDecrypt(b *testing.B) {
	ivValue := uint32(42)
	iv2 := ivValue ^ IVSeedFormatC

	var nonce [8]byte
	binary.LittleEndian.PutUint32(nonce[0:4], iv2)
	binary.LittleEndian.PutUint32(nonce[4:8], ivValue)

	cleartext := make([]byte, 0x170)
	binary.LittleEndian.PutUint32(cleartext[0:4], magicNumber)

	encrypted := make([]byte, len(cleartext))
	salsa20.XORKeyStream(encrypted, cleartext, nonce[:], getSalsa20Key())
	binary.LittleEndian.PutUint32(encrypted[ivOffset:ivOffset+4], ivValue)

	b.ResetTimer()

	for b.Loop() {
		_, _ = decrypt(encrypted, IVSeedFormatC)
	}
}

func FuzzDecrypt(f *testing.F) {
	f.Add(make([]byte, 0x170))

	ivValue := uint32(42)
	iv2 := ivValue ^ IVSeedFormatC

	var nonce [8]byte
	binary.LittleEndian.PutUint32(nonce[0:4], iv2)
	binary.LittleEndian.PutUint32(nonce[4:8], ivValue)

	cleartext := make([]byte, 0x170)
	binary.LittleEndian.PutUint32(cleartext[0:4], magicNumber)

	encrypted := make([]byte, len(cleartext))
	salsa20.XORKeyStream(encrypted, cleartext, nonce[:], getSalsa20Key())
	binary.LittleEndian.PutUint32(encrypted[ivOffset:ivOffset+4], ivValue)
	f.Add(encrypted)

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = decrypt(data, IVSeedFormatC)
	})
}
