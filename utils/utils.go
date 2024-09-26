package utils

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"math/big"
)

func Prepend[T any](x []T, y T) []T {
	var temp T
	x = append(x, temp)
	copy(x[1:], x)
	x[0] = y
	return x
}

func MakeMap[K any]() map[string]K {
	return make(map[string]K)
}

func GenerateRandomId(length int) (string, error) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func RandomFloat64() (float64, error) {
	// Generate 8 random bytes
	var buf [8]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return 0, err
	}

	// Convert bytes to uint64 and then to a float64
	randUint := binary.BigEndian.Uint64(buf[:])

	// Normalize the result to be in the range [0, 1)
	return float64(randUint) / float64(math.MaxUint64), nil
}
