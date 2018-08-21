package b58

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func hex2bytes(hexstring string) (b []byte) {
	b, _ = hex.DecodeString(hexstring)
	return b
}

func TestBase58(t *testing.T) {
	var b58Vectors = []struct {
		bytes   []byte
		encoded string
	}{
		{hex2bytes("4e19"), "6wi"},
		{hex2bytes("3ab7"), "5UA"},
		{hex2bytes("ae0ddc9b"), "5T3W5p"},
		{hex2bytes("65e0b4c9"), "3c3E6L"},
		{hex2bytes("25793686e9f25b6b"), "7GYJp3ZThFG"},
		{hex2bytes("94b9ac084a0d65f5"), "RspedB5CMo2"},
	}

	/* Test base-58 encoding */
	for i := 0; i < len(b58Vectors); i++ {
		got := Encode(b58Vectors[i].bytes)
		if got != b58Vectors[i].encoded {
			t.Fatalf("Encode(%v): got %s, expected %s", b58Vectors[i].bytes, got, b58Vectors[i].encoded)
		}
	}
	t.Log("success Encode() on valid vectors")

	/* Test base-58 decoding */
	for i := 0; i < len(b58Vectors); i++ {
		got, err := Decode(b58Vectors[i].encoded)
		if err != nil {
			t.Fatalf("Decode(%s): got error %v, expected %v", b58Vectors[i].encoded, err, b58Vectors[i].bytes)
		}
		if bytes.Compare(got, b58Vectors[i].bytes) != 0 {
			t.Fatalf("Decode(%s): got %v, expected %v", b58Vectors[i].encoded, got, b58Vectors[i].bytes)
		}
	}
	t.Log("success Decode() on valid vectors")

	/* Test base-58 decoding of invalid strings */
	b58InvalidVectors := []string{
		"5T3IW5p", // Invalid character I
		"6Owi",    // Invalid character O
	}

	for i := 0; i < len(b58InvalidVectors); i++ {
		got, err := Decode(b58InvalidVectors[i])
		if err == nil {
			t.Fatalf("Decode(%s): got %v, expected error", b58InvalidVectors[i], got)
		}
		t.Logf("Decode(%s): got expected err %v", b58InvalidVectors[i], err)
	}
	t.Log("success Decode() on invalid vectors")
}
