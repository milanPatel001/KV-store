package utils

import (
	"crypto/rand"
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
