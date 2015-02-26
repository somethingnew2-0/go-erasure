// Erasure coding is similar to RAID based parity encoding, but is more generalized
// and powerful. When defining an erasure code, you specify a K and M variable.
// M is the number of shards you wish to encode and K is the number shards it
// takes to recreate your original data. Hence K must be less than M and usually
// not equal (as that would be a pointless encoding). The real magic with erasure
// coding is that fact that ANY K of the M shards can recreate the original data.
// For example, a erasure coding scheme of K=8 and M=12 means any four of the
// encoded shards can be lost while the original data can still be constructed
// from the valid remaining eight shards.
package erasure

// #cgo CFLAGS: -Wall -std=gnu99
// #include "types.h"
// #include "erasure_code.h"
import "C"

import (
	"sync"
)

// Manages state of the erasure coding scheme and its values should be
// considered read-only.
type Code struct {
	M            int
	K            int
	ShardLength  int
	EncodeMatrix []byte
	galoisTables []byte
	decode       *decodeNode
}

type decodeNode struct {
	children     []*decodeNode
	mutex        *sync.Mutex
	galoisTables []byte
	decodeIndex  []byte
}

func (c *Code) getDecode(errList []byte) *decodeNode {
	node := c.decode.getDecode(errList, 0, byte(c.M))

	node.mutex.Lock()
	defer node.mutex.Unlock()
	if node.galoisTables == nil || node.decodeIndex == nil {
		node.galoisTables = make([]byte, c.K*c.M*32)
		node.decodeIndex = make([]byte, c.K)

		decodeMatrix := make([]byte, c.M*c.K)
		srcInErr := make([]byte, c.M)
		nErrs := len(errList)
		nSrcErrs := 0
		for _, err := range errList {
			srcInErr[err] = 1
			if int(err) < c.K {
				nSrcErrs++
			}
		}

		C.gf_gen_decode_matrix((*C.uchar)(&c.EncodeMatrix[0]), (*C.uchar)(&decodeMatrix[0]), (*C.uchar)(&node.decodeIndex[0]), (*C.uchar)(&errList[0]), (*C.uchar)(&srcInErr[0]), C.int(nErrs), C.int(nSrcErrs), C.int(c.K), C.int(c.M))

		C.ec_init_tables(C.int(c.K), C.int(nErrs), (*C.uchar)(&decodeMatrix[0]), (*C.uchar)(&node.galoisTables[0]))
	}

	return node
}

func (n *decodeNode) getDecode(errList []byte, parent, m byte) *decodeNode {
	n.mutex.Lock()
	node := n.children[errList[0]-parent]
	if node == nil {
		node = &decodeNode{
			children: make([]*decodeNode, m-errList[0]),
			mutex:    &sync.Mutex{},
		}
		n.children[errList[0]-parent] = node
	}
	n.mutex.Unlock()
	if len(errList) > 1 {
		return node.getDecode(errList[1:], errList[0]+1, m)
	}
	return node
}

// Constructor for creating a new erasure coding scheme. M is the total
// number of shards output by the encoding.  K is the number of shards
// that can recreate any data that was encoded.  Size is the size of the
// byte array to encode.  It should be divisible by K as each shard
// will be Size / K in length.  The maximum value for K and M is 127.
func NewCode(m int, k int, size int) *Code {
	if m <= 0 || k <= 0 || k >= m || k > 127 || m > 127 || size < 0 {
		panic("Invalid erasure code params")
	}
	if size%k != 0 {
		panic("Size to encode is not divisable by k and therefore cannot be encoded into shards")
	}

	encodeMatrix := make([]byte, m*k)
	galoisTables := make([]byte, k*(m-k)*32)

	if k > 5 {
		C.gf_gen_cauchy1_matrix((*C.uchar)(&encodeMatrix[0]), C.int(m), C.int(k))
	} else {
		C.gf_gen_rs_matrix((*C.uchar)(&encodeMatrix[0]), C.int(m), C.int(k))
	}

	C.ec_init_tables(C.int(k), C.int(m-k), (*C.uchar)(&encodeMatrix[k*k]), (*C.uchar)(&galoisTables[0]))
	return &Code{
		M:            m,
		K:            k,
		ShardLength:  size / k,
		EncodeMatrix: encodeMatrix,
		galoisTables: galoisTables,
		decode: &decodeNode{
			children: make([]*decodeNode, m),
			mutex:    &sync.Mutex{},
		},
	}
}

// The data buffer to encode must be of the length Size given in the constructor.
// The returned encoded buffer is (M-K)*Shard length, since the first Size bytes
// of the encoded data is just the original data due to the identity matrix.
func (c *Code) Encode(data []byte) []byte {
	if len(data) != c.K*c.ShardLength {
		panic("Data to encode is not the proper size")
	}
	// Since the first k rows of the encode matrix is actually the identity matrix
	// we only need to encode the last m-k shards of the data and append
	// them to the original data
	encoded := make([]byte, (c.M-c.K)*(c.ShardLength))
	C.ec_encode_data(C.int(c.ShardLength), C.int(c.K), C.int(c.M-c.K), (*C.uchar)(&c.galoisTables[0]), (*C.uchar)(&data[0]), (*C.uchar)(&encoded[0]))
	// return append(data, encoded...)
	return encoded
}

// Data buffer to decode must be of the (M/K)*Size given in the constructor.
// The error list must contain M-K values, corresponding to the shards
// with errors (eg. [0, 2, 4, 6]).
// The returned decoded data is the orignal data of length Size
func (c *Code) Decode(encoded []byte, errList []byte) []byte {
	if len(encoded) != c.M*c.ShardLength {
		panic("Data to decode is not the proper size")
	}
	if len(errList) > c.M-c.K {
		panic("Too many errors, cannot decode")
	}
	recovered := []byte{}
	if len(errList) == 0 {
		recovered = append(recovered, encoded[:c.K*c.ShardLength]...)
	} else {
		node := c.getDecode(errList)

		for i := 0; i < c.K; i++ {
			recovered = append(recovered, encoded[(int(node.decodeIndex[i])*c.ShardLength):int(node.decodeIndex[i]+1)*c.ShardLength]...)
		}

		decoded := make([]byte, c.M*c.ShardLength)
		C.ec_encode_data(C.int(c.ShardLength), C.int(c.K), C.int(c.M), (*C.uchar)(&node.galoisTables[0]), (*C.uchar)(&recovered[0]), (*C.uchar)(&decoded[0]))

		copy(recovered, encoded)

		for i, err := range errList {
			if int(err) < c.K {
				copy(recovered[int(err)*c.ShardLength:int(err+1)*c.ShardLength], decoded[i*c.ShardLength:(i+1)*c.ShardLength])
			}
		}
	}

	return recovered
}
