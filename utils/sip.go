package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
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

// Error Error
type Error struct {
	err    error
	params []interface{}
}

func (err *Error) Error() string {
	if err == nil {
		return "<nil>"
	}
	str := fmt.Sprint(err.params...)
	if err.err != nil {
		str += fmt.Sprintf(" err:%s", err.err.Error())
	}
	return str
}

// NewError NewError
func NewError(err error, params ...interface{}) error {
	return &Error{err, params}
}
func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}

// ResolveSelfIP ResolveSelfIP
func ResolveSelfIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip, nil
		}
	}
	return nil, errors.New("server not connected to any network")
}
