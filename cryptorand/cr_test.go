package cryptorand

import "testing"

func TestMustInt64(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.FailNow()
		}
	}()
	for i := 0; i < 64; i++ {
		if MustInt64() == MustInt64() {
			t.Log("MustInt64 must never return the same value twice")
		}
	}
}

func TestInt64(t *testing.T) {
	for i := 0; i < 64; i++ {
		x, err1 := Int64()
		y, err2 := Int64()
		if err1 != nil {
			t.Logf("First Int64 call returned an error at iteration #%d", i)
		}
		if err2 != nil {
			t.Logf("Second Int64 call returned an error at iteration #%d", i)
		}
		if x == y {
			t.Log("Int64 must never return the same value twice")
		}
	}
}

func TestMustFloat64(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.FailNow()
		}
	}()
	for i := 0; i < 64; i++ {
		if MustFloat64() == MustFloat64() {
			t.Log("MustFloat64 must never return the same value twice")
		}
	}
}

func TestFloat64(t *testing.T) {
	for i := 0; i < 64; i++ {
		x, err1 := Float64()
		y, err2 := Float64()

		if err1 != nil {
			t.Logf("First Float64 call returned an error at iteration #%d", i)
		}
		if err2 != nil {
			t.Logf("Second Float64 call returned an error at iteration #%d", i)
		}
		if x == y {
			t.Log("Float64 must never return the same value twice")
		}
	}
}
