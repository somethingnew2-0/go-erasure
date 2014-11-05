package main

import (
	"erasure"

	"log"
)

func main() {
	m := 12
	k := 8
	size := 8 * 16
	// erasure.Hello()
	code := erasure.NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = 0x62
	}

	encoded := code.Encode(source)

	log.Printf("Encoded: %x\n", encoded)
	srcErrList := []int8{0, 2, 3, 4}

	recovered := code.Decode(encoded, srcErrList)
	log.Printf("Recovered: %x\n", recovered)
}
