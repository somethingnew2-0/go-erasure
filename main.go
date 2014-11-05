package main

import (
	"erasure"

	"log"
	"math/rand"
)

func main() {
	m := 12
	k := 8
	size := 8 * 16
	// erasure.Hello()
	code := erasure.NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	log.Printf("Source: %x\n", source)

	encoded := code.Encode(source)

	log.Printf("Encoded: %x\n", encoded)
	srcErrList := []int8{0, 2, 3, 4}

	for _, err := range srcErrList {
		for i := 0; i < code.VectorLength; i++ {
			source[int(err)*code.VectorLength+i] = 0x62
		}
	}

	log.Printf("Source Corrupted: %x\n", source)

	recovered := code.Decode(append(source, encoded...), srcErrList)
	log.Printf("Recovered: %x\n", recovered)

}
