package utils

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func RandNum16String(n int) string {
	numbers16 := "0123456789abcdef"
	return randStringBySoure(numbers16, n)
}

func RandNumString(n int) string {
	numbers := "0123456789"
	return randStringBySoure(numbers, n)
}

func RandString(n int) string {
	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return randStringBySoure(letterBytes, n)
}

// https://github.com/kpbird/golang_random_string
func randStringBySoure(src string, n int) string {
	randomness := make([]byte, n)

	rand.Seed(time.Now().UnixNano())
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}

	l := len(src)

	// fill output
	output := make([]byte, n)
	for pos := range output {
		random := randomness[pos]
		randomPos := random % uint8(l)
		output[pos] = src[randomPos]
	}

	return string(output)
}

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}
