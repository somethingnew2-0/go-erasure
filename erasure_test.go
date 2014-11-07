package erasure

import (
	"bytes"
	"math/rand"
	"testing"
)

func corrupt(source, errList []byte, vectorLength int) []byte {
	corrupted := make([]byte, len(source))
	copy(corrupted, source)
	for _, err := range errList {
		for i := 0; i < vectorLength; i++ {
			corrupted[int(err)*vectorLength+i] = 0x00
		}
	}
	return corrupted
}

func TestErasure_12_8(t *testing.T) {
	m := 12
	k := 8
	vectorLength := 16
	size := k * vectorLength

	code := NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	encoded := code.Encode(source)

	errList := []byte{0, 2, 3, 4}

	corrupted := corrupt(append(source, encoded...), errList, vectorLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 4 errors")
	}
}

func TestErasure_16_8(t *testing.T) {
	m := 16
	k := 8
	vectorLength := 16
	size := k * vectorLength

	code := NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	encoded := code.Encode(source)

	errList := []byte{0, 1, 2, 3, 4, 5, 6, 7}

	corrupted := corrupt(append(source, encoded...), errList, vectorLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 8 errors")
	}
}

func TestErasure_17_8(t *testing.T) {
	m := 17
	k := 8
	vectorLength := 16
	size := k * vectorLength

	code := NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	encoded := code.Encode(source)

	errList := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}

	corrupted := corrupt(append(source, encoded...), errList, vectorLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 4 errors")
	}
}

func TestErasure_9_5(t *testing.T) {
	m := 9
	k := 5
	vectorLength := 16
	size := k * vectorLength

	code := NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	encoded := code.Encode(source)

	errList := []byte{0, 2, 3, 4}

	corrupted := corrupt(append(source, encoded...), errList, vectorLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 4 errors")
	}
}
