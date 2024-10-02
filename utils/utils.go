package utils

import (
	"bytes"
	"cmp"
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"math"
	"math/big"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var CommandsWithRequiredArgs []string = []string{"SET", "DEL", "GET", "NUM"}

func StoreCacheGobEncoded[K string | int, V any](fileName string, cache map[K]V) error {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(cache); err != nil {
		fmt.Println(err)
		return err
	}

	if err := CreatDir("./temp"); err != nil {
		return err
	}

	completeFileName := fmt.Sprintf("./temp/%s.gob", fileName)

	file, err := os.Create(completeFileName)

	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	if _, err = file.Write(buf.Bytes()); err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}

func DecodeGobFile[K string | int, V any](fileName string) (map[K]V, error) {
	var m map[K]V
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)

	completeFileName := fmt.Sprintf("./temp/%s.gob", fileName)

	file, err := os.Open(completeFileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	_, err = buf.ReadFrom(file)
	if err != nil {
		fmt.Println("Error reading from file:", err)
		return nil, err
	}

	if err = dec.Decode(&m); err != nil {
		fmt.Println("Error decoding:", err)
		return nil, err
	}

	return m, nil
}

func SerializeOutput(command string, commandOutput string) string {
	return fmt.Sprintf("%v\r\n%v\r\n", command, commandOutput)
}

func DeserializeInput(str string) (string, []string, error) {
	in := strings.Split(str, "\r\n")

	if len(in) == 0 {
		return "", nil, fmt.Errorf(">> Nothing Entered !!!")
	}

	command := strings.ToUpper(in[0])
	args := in[1 : len(in)-1]

	return command, args, nil
}

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

func SetMaxValue[T cmp.Ordered]() T {

	typ := reflect.TypeFor[T]()

	maxs := [...]any{
		reflect.Int:   math.MaxInt,
		reflect.Int8:  math.MaxInt8,
		reflect.Int16: math.MaxInt16,
		reflect.Int32: math.MaxInt32,
		reflect.Int64: math.MaxInt64,

		reflect.Uint:   uint(math.MaxUint),
		reflect.Uint8:  math.MaxUint8,
		reflect.Uint16: math.MaxUint16,
		reflect.Uint32: math.MaxUint32,
		reflect.Uint64: uint64(math.MaxUint64),

		reflect.Float32: math.MaxFloat32,
		reflect.Float64: math.MaxFloat64,

		reflect.String: "INF",
	}

	v := maxs[typ.Kind()]
	val := reflect.ValueOf(v).Convert(typ)
	return val.Interface().(T)
}

func CreatDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist, creating: %s\n", dirPath)

		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}
	} else {
		fmt.Println("Directory already exists!")
	}

	return nil
}

func LocateGobFile(dirPath, pattern string) (string, error) {

	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return "", err
	}

	fullPattern := fmt.Sprintf(`^%v\.gob$`, pattern)
	matcher := regexp.MustCompile(fullPattern)

	for _, file := range files {
		if !file.IsDir() && matcher.MatchString(file.Name()) {
			return file.Name(), nil
		}
	}

	return "", fmt.Errorf("No gob files exist !!!")
}
