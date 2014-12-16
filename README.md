go-erasure
========

Go bindings for erasure coding (Reed-Solomon coding).

Erasure coding is similar to RAID based parity encoding, but is more generalized and powerful.  When defining an erasure code, you specify a `k` and `m` variable. `m` is the number of shards you wish to encode and `k` is the number shards it takes to recreate your original data.  Hence `k` must be less than `m` and usually not equal (as that would be a pointless encoding). The real magic with erasure coding is that fact that ANY `k` of the `m` shards can recreate the original data.  For example, a erasure coding scheme of `k=8` and `m=12` means any four of the encoded shards can be lost while the original data can still be constructed from the valid remaining eight shards.

This library is aimed at simplicity and performance.  It only has three methods including a constructor which are all thread-safe! Internally it uses cgo to utilize a complex C library.  For more indepth look into this library be sure to check out the [Intel® Storage Acceleration Library](https://01.org/intel%C2%AE-storage-acceleration-library-open-source-version) and especially their corresponding [video](http://www.intel.com/content/www/us/en/storage/erasure-code-isa-l-solution-video.html).  One feature it does add is an optimization for decoding.  Since there are `m choose k` possible inverse matrices for decoding, this library caches them (via lazy-loading) so as reduce the amount of time decoding.  It does so by utilizing a [trie](http://en.wikipedia.org/wiki/Trie) where the sorted error list of shards is the key to the trie and the corresponding decode matrix is the value.

I hope you find it useful and pull requests are welcome!

## Usage
See the [GoDoc](https://godoc.org/github.com/somethingnew2-0/go-erasure) for an API reference

### Encode and decode random data

```go
package main

import (
  "bytes"
  "log"
  "math/rand"
  
  "github.com/somethingnew2-0/go-erasure"
)

func corrupt(source, errList []byte, shardLength int) []byte {
	corrupted := make([]byte, len(source))
	copy(corrupted, source)
	for _, err := range errList {
		for i := 0; i < shardLength; i++ {
			corrupted[int(err)*shardLength+i] = 0x00
		}
	}
	return corrupted
}

func main() {
	m := 12
	k := 8
	shardLength := 16 // Length of a shard
	size := k * shardLength // Length of the data blob to encode

	code := erasure.NewCode(m, k, size)

	source := make([]byte, size)
	for i := range source {
		source[i] = byte(rand.Int63() & 0xff) //0x62
	}

	encoded := code.Encode(source)

	errList := []byte{0, 2, 3, 4}

	corrupted := corrupt(append(source, encoded...), errList, shardLength)

	recovered := code.Decode(corrupted, errList)

	if !bytes.Equal(source, recovered) {
		log.Fatal("Source was not sucessfully recovered with 4 errors")
	}
}
```


## Development

To start run `source install.sh` or more simply `. install.sh` to setup the git hooks and GOPATH for this project.

Run `go test` or `go test -bench .` to test the unit tests and benchmark tests.
