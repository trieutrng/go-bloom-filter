package bloomfilter

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"hash"
	"hash/fnv"
	"math"
)

const FALSE_POSITIVE_RATE = 0.000001

type BloomFilter interface {
	Add(value interface{})
	MightContain(value interface{}) bool
}

type bloomFilterImpl struct {
	filter    []bool
	hashChain []hash.Hash
	size      int64
}

/*
	 params
		- cap: expected capacity of data
*/
func New(cap int64) BloomFilter {
	size := getFilterSize(cap)
	return &bloomFilterImpl{
		filter:    make([]bool, size),
		hashChain: getHashChain(),
		size:      size,
	}
}

func (bf *bloomFilterImpl) Add(value interface{}) {
	for _, hashFunc := range bf.hashChain {
		index := hashing(value, hashFunc) % bf.size
		bf.filter[index] = true
	}
}

func (bf *bloomFilterImpl) MightContain(value interface{}) bool {
	for _, hashFunc := range bf.hashChain {
		index := hashing(value, hashFunc)
		if !bf.filter[index] {
			return false
		}
	}
	return true
}

// this fomular is referenced from https://redis.io/docs/latest/develop/data-types/probabilistic/bloom-filter/
func getFilterSize(cap int64) int64 {
	bitsPerItem := -math.Log(FALSE_POSITIVE_RATE) / math.Log(2)
	return int64(bitsPerItem) * cap
}

func hashing(value interface{}, hashFunc hash.Hash) int64 {
	hashFunc.Write([]byte(encoding(value)))
	bits := hashFunc.Sum(nil)
	buffer := bytes.NewBuffer(bits)
	result, _ := binary.ReadVarint(buffer)
	hashFunc.Reset()
	return result
}

func getHashChain() []hash.Hash {
	return []hash.Hash{
		fnv.New32(),
		fnv.New32(),
		fnv.New32(),
	}
}

func encoding(value interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(value)
	return buf.Bytes()
}
