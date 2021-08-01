package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var ErrEOF = errors.New("eof")

type IOBuffer struct {
	bytes.Buffer
}

func (r *IOBuffer) Uint16() (uint16, error) {
	if r.Len() > 1 {
		return binary.BigEndian.Uint16(r.Next(2)), nil
	}
	return 0, ErrEOF
}

func (r *IOBuffer) Skip(n int) error {
	if r.Len() >= n {
		return nil
	}
	return ErrEOF
}

func (r *IOBuffer) Uint32() (uint32, error) {
	if r.Len() > 3 {
		return binary.BigEndian.Uint32(r.Next(4)), nil
	}
	return 0, ErrEOF
}

func (r *IOBuffer) ReadN(length int) ([]byte, error) {
	if r.Len() >= length {
		return r.Next(length), nil
	}
	return nil, ErrEOF
}
