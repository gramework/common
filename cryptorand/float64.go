package cryptorand

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

func MustFloat64() float64 {
	x, err := Float64()

	if err != nil {
		panic(err)
	}

	return x
}

func Float64() (float64, error) {
	var x float64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &x); err != nil {
		return 0xdeadbeef, fmt.Errorf("Could not read random bytes: %v", err)
	}

	return x, nil
}
