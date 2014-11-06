package erasure

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestErasure_12_8(t *testing.T) {
	m := 12
	k := 8
	size := k * 16

	code := NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	t.Logf("Source: %x\n", source)

	encoded := code.Encode(source)

	t.Logf("Encoded: %x\n", encoded)
	srcErrList := []int8{0, 2, 3, 4}

	corrupted := make([]byte, size)
	copy(corrupted, source)
	for _, err := range srcErrList {
		for i := 0; i < code.VectorLength; i++ {
			corrupted[int(err)*code.VectorLength+i] = 0x62
		}
	}

	t.Logf("Source Corrupted: %x\n", corrupted)

	recovered := code.Decode(append(corrupted, encoded...), srcErrList)
	t.Logf("Recovered: %x\n", recovered)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 4 errors")
	}
}

// func TestErasure_9_5(t *testing.T) {
// 	m := 9
// 	k := 5
// 	size := k * 16

// 	code := NewCode(m, k, size)

// 	source := make([]byte, size)
// 	for i := range source {
// 		source[i] = byte(rand.Int63() & 0xff) //0x62
// 	}

// t.Logf("Source: %x\n", source)

// encode := code.Encode(source)

// t.Logf("Encoded: %x\n", encoded)
// srcErrList := []int8{0, 2, 3, 4}

// corrupted := make([]byte, size)
// copy(corrupted, source)
// for _, err := range srcErrList {
// 	for i := 0; i < code.VectorLength; i++ {
// 		corrupted[int(err)*code.VectorLength+i] = 0x62
// 	}
// }

// t.Logf("Source Corrupted: %x\n", corrupted)

// recovered := code.Decode(append(corrupted, encoded...), srcErrList)
// t.Logf("Recovered: %x\n", recovered)

// if !bytes.Equal(source, recovered) {
// 	t.Error("Source was not sucessfully recovered with 4 errors")
// }
// }
