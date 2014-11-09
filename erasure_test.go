package erasure

import (
	"bytes"
	"math/rand"
	"runtime"
	"sort"
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

func randomErrorList(m, numberOfErrs int) []byte {
	set := make(map[int]bool, m)
	errListInts := make([]int, numberOfErrs)
	for i := 0; i < numberOfErrs; i++ {
		err := rand.Intn(m)
		for set[err] {
			err = rand.Intn(m)
		}
		set[err] = true
		errListInts[i] = err
	}

	sort.Ints(errListInts)

	errList := make([]byte, numberOfErrs)
	for i, err := range errListInts {
		errList[i] = byte(err)
	}

	return errList
}

func TestBasicErasure_12_8(t *testing.T) {
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

func TestBasicErasure_16_8(t *testing.T) {
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

func TestBasicErasure_20_8(t *testing.T) {
	m := 20
	k := 8
	vectorLength := 16
	size := k * vectorLength

	code := NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	encoded := code.Encode(source)

	errList := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 16, 17}

	corrupted := corrupt(append(source, encoded...), errList, vectorLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 4 errors")
	}
}

func TestBasicErasure_9_5(t *testing.T) {
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

func TestRandomErasure_12_8(t *testing.T) {
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

	errList := randomErrorList(m, rand.Intn(m-k))

	corrupted := corrupt(append(source, encoded...), errList, vectorLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 4 errors")
	}
}

func TestRandomErasure_16_8(t *testing.T) {
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

	errList := randomErrorList(m, rand.Intn(m-k))

	corrupted := corrupt(append(source, encoded...), errList, vectorLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 8 errors")
	}
}

func TestRandomErasure_20_8(t *testing.T) {
	m := 20
	k := 8
	vectorLength := 16
	size := k * vectorLength

	code := NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	encoded := code.Encode(source)

	errList := randomErrorList(m, rand.Intn(m-k))

	corrupted := corrupt(append(source, encoded...), errList, vectorLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 4 errors")
	}
}

func TestRandomErasure_9_5(t *testing.T) {
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

	errList := randomErrorList(m, rand.Intn(m-k))

	corrupted := corrupt(append(source, encoded...), errList, vectorLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		t.Error("Source was not sucessfully recovered with 4 errors")
	}
}

func BenchmarkBasicEncode_12_8(b *testing.B) {
	m := 12
	k := 8
	vectorLength := 16
	size := k * vectorLength

	code := NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	b.SetBytes(int64(size))
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			code.Encode(source)
		}
	})
}

func BenchmarkBasicDecode_12_8(b *testing.B) {
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

	b.SetBytes(int64(size))
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		errList := []byte{0, 2, 3, 4}

		corrupted := corrupt(append(source, encoded...), errList, vectorLength)

		for pb.Next() {
			recovered := code.Decode(corrupted, errList)

			if !bytes.Equal(source, recovered) {
				b.Error("Source was not sucessfully recovered with 4 errors")
			}
		}
	})
}

func BenchmarkRandomDecode_12_8(b *testing.B) {
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

	b.SetBytes(int64(size))
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			errList := randomErrorList(m, rand.Intn(m-k))

			corrupted := corrupt(append(source, encoded...), errList, vectorLength)

			recovered := code.Decode(corrupted, errList)

			if !bytes.Equal(source, recovered) {
				b.Error("Source was not sucessfully recovered with 4 errors")
			}
		}
	})
}
