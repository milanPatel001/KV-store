package utils

import (
	"fmt"
	"math"

	"github.com/cespare/xxhash/v2"
)

type BloomFilter struct {
	Arr               []byte
	EstimatedCapacity uint
	HashFuncNum       uint8
}

func CreateBloomFilter(capacity uint, errorRate float64) *BloomFilter {

	requiredBits := math.Floor(-float64(capacity) * math.Log(errorRate) / (math.Ln2 * math.Ln2))
	byteConv := GetBloomFilterArrSize(uint(requiredBits))

	hashfuncNum := math.Floor(requiredBits / float64((capacity)) * math.Ln2)

	bl := BloomFilter{EstimatedCapacity: capacity, HashFuncNum: uint8(hashfuncNum)}
	bl.Arr = make([]byte, int(byteConv))

	return &bl
}

func (bl *BloomFilter) Set(key string) error {
	if key == "" {
		return fmt.Errorf("Key can't be empty !!!")
	}

	for i := 0; i < int(bl.HashFuncNum); i++ {
		hash := xxhash.Sum64([]byte(key)) + uint64(i)*0x9E3779B97F4A7C15

		bitIndex := hash % uint64(len(bl.Arr)*8)

		bl.SetBit(int(bitIndex))
	}

	return nil
}

func (bl *BloomFilter) DoesExist(key string) bool {

	for i := 0; i < int(bl.HashFuncNum); i++ {
		hash := xxhash.Sum64([]byte(key)) + uint64(i)*0x9E3779B97F4A7C15

		bitIndex := hash % uint64(len(bl.Arr)*8)

		if bl.GetBit(int(bitIndex)) == false {
			return false
		}
	}

	return true
}

// bitIndex starting from 0
func (bl *BloomFilter) SetBit(bitIndex int) error {
	if len(bl.Arr)*8 <= bitIndex {
		return fmt.Errorf("Index out of bounds !!!")
	}

	arrIndex := bitIndex / 8
	shift := (8 - bitIndex%8 - 1)

	bl.Arr[arrIndex] = (1 << shift) | bl.Arr[arrIndex]

	return nil
}

func (bl *BloomFilter) GetBit(bitIndex int) bool {
	if len(bl.Arr)*8 <= int(bitIndex) {
		return false
	}

	arrIndex := bitIndex / 8
	shift := (8 - bitIndex%8 - 1)

	return ((1 << shift) & bl.Arr[arrIndex]) != 0
}

// MaxCapacity -> number of bits for your bloom filter
func GetBloomFilterArrSize(maxCapacity uint) int {

	if maxCapacity < 16 {
		maxCapacity = 16
	}

	byteConv := math.Ceil(float64(maxCapacity) / 8)

	return int(byteConv)

}

// Decide on creating new sub-stack of filter if whole arr is filled
