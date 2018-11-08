package cryptorand

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

func MustInt64() int64 {
	x, err := Int64()

	if err != nil {
		panic(err)
	}

	return x
}

func Int64() (int64, error) {
	var x int64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &x); err != nil {
		return 0xdeadbeef, fmt.Errorf("Could not read random bytes: %v", err)
	}

	return x, nil
}

func MustUInt64() uint64 {
	x, err := UInt64()

	if err != nil {
		panic(err)
	}

	return x
}

func UInt64() (uint64, error) {
	var x uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &x); err != nil {
		return 0xdeadbeef, fmt.Errorf("Could not read random bytes: %v", err)
	}

	return x, nil
}
