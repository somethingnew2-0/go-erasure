package erasure

/*
#cgo CFLAGS: -Wall -std=gnu99
#include "types.h"
#include "erasure_code.h"
*/
import "C"

import (
	"log"
)

type Code struct {
	M            int
	K            int
	VectorLength int
	EncodeMatrix []byte
	galoisTables []byte
}

func NewCode(m int, k int, size int) *Code {
	if m <= 0 || k <= 0 || k >= m || k > 127 || m > 127 || size < 0 {
		log.Fatal("Invalid erasure code params")
	}
	if size%k != 0 {
		log.Fatal("Size to encode is not divisable by k and therefore cannot be enocded in vector chunks")
	}

	encodeMatrix := make([]byte, m*k)
	galoisTables := make([]byte, k*(m-k)*32)

	if k > 5 {
		C.gf_gen_cauchy1_matrix((*C.uchar)(&encodeMatrix[0]), C.int(m), C.int(k))
	} else {
		C.gf_gen_rs_matrix((*C.uchar)(&encodeMatrix[0]), C.int(m), C.int(k))
	}

	return &Code{
		M:            m,
		K:            k,
		VectorLength: size / k,
		EncodeMatrix: encodeMatrix,
		galoisTables: galoisTables,
	}
}

// Data buffer to encode must be of the k*vector given in the constructor
// The returned encoded buffer is (m-k)*vector, since the first k*vector of the
// encoded data is just the original data due to the identity matrix
func (c *Code) Encode(data []byte) []byte {
	if len(data) != c.K*c.VectorLength {
		log.Fatal("Data to encode is not the proper size")
	}
	// Since the first k row of the encode matrix is actually the identity matrix
	// we only need to encode the last m-k vectors of the matrix and append
	// them to the original data
	encoded := make([]byte, (c.M-c.K)*(c.VectorLength))
	C.ec_init_tables(C.int(c.K), C.int(c.M-c.K), (*C.uchar)(&c.EncodeMatrix[c.K*c.K]), (*C.uchar)(&c.galoisTables[0]))
	C.ec_encode_data(C.int(c.VectorLength), C.int(c.K), C.int(c.M-c.K), (*C.uchar)(&c.galoisTables[0]), (*C.uchar)(&data[0]), (*C.uchar)(&encoded[0]))

	// return append(data, encoded...)
	return encoded
}

// Data buffer to decode must be of the m*vector given in the constructor
// The source error list must contain m-k values, corresponding to the vectors with errors
// The returned decoded data is k*vector
func (c *Code) Decode(encoded []byte, srcErrList []byte) []byte {
	if len(encoded) != c.M*c.VectorLength {
		log.Fatal("Data to decode is not the proper size")
	}
	if len(srcErrList) > c.M-c.K {
		log.Fatal("Too many errors, cannot decode")
	}
	decodeMatrix := make([]byte, c.M*c.K)
	decodeIndex := make([]byte, c.K)
	srcInErr := make([]byte, c.M)
	nErrs := len(srcErrList)
	nSrcErrs := 0
	for _, err := range srcErrList {
		srcInErr[err] = 1
		if int(err) < c.K {
			nSrcErrs++
		}
	}

	C.gf_gen_decode_matrix((*C.uchar)(&c.EncodeMatrix[0]), (*C.uchar)(&decodeMatrix[0]), (*C.uchar)(&decodeIndex[0]), (*C.uchar)(&srcErrList[0]), (*C.uchar)(&srcInErr[0]), C.int(nErrs), C.int(nSrcErrs), C.int(c.K), C.int(c.M))

	C.ec_init_tables(C.int(c.K), C.int(nErrs), (*C.uchar)(&decodeMatrix[0]), (*C.uchar)(&c.galoisTables[0]))

	recovered := []byte{}
	for i := 0; i < c.K; i++ {
		recovered = append(recovered, encoded[(int(decodeIndex[i])*c.VectorLength):int(decodeIndex[i]+1)*c.VectorLength]...)
	}

	decoded := make([]byte, c.M*c.VectorLength)
	C.ec_encode_data(C.int(c.VectorLength), C.int(c.K), C.int(c.M), (*C.uchar)(&c.galoisTables[0]), (*C.uchar)(&recovered[0]), (*C.uchar)(&decoded[0]))

	copy(recovered, encoded)

	for i, err := range srcErrList {
		copy(recovered[int(err)*c.VectorLength:int(err+1)*c.VectorLength], decoded[i*c.VectorLength:(i+1)*c.VectorLength])
	}

	return recovered
}
