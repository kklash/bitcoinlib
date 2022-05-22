package blockheader

import (
	"math/big"
	"strings"
	"testing"
)

func TestTargetNBits(t *testing.T) {
	type Fixture struct {
		nBits       uint32
		targetNBits string
	}

	fixtures := []*Fixture{
		{0x0, "0x0"},
		{0x00123456, "0x0"},
		{0x00923456, "0x0"},
		{0x01003456, "0x0"},
		{0x01123456, "0x12"},
		{0x01803456, "0x0"},
		{0x01fedcba, "0x0"},
		{0x02000056, "0x0"},
		{0x02008000, "0x80"},
		{0x02123456, "0x1234"},
		{0x02800056, "0x0"},
		{0x03000000, "0x0"},
		{0x03123456, "0x123456"},
		{0x03800000, "0x0"},
		{0x04000000, "0x0"},
		{0x04123456, "0x12345600"},
		{0x04800000, "0x0"},
		{0x04923456, "0x0"},
		{0x05009234, "0x92340000"},
		{0x1d00ffff, "0xFFFF0000000000000000000000000000000000000000000000000000"},
		{0x20123456, "0x1234560000000000000000000000000000000000000000000000000000000000"},
		{0xff123456, ""},
	}

	for _, fixture := range fixtures {
		expected, _ := new(big.Int).SetString(fixture.targetNBits, 0)

		if expected == nil {
			defer func() {
				panicValue := recover()
				if panicValue == nil {
					return
				}

				panicString, ok := panicValue.(string)
				if !ok || fixture.targetNBits != "" {
					t.Errorf("unexpected panic value: %v", panicValue)
				} else if !strings.Contains(panicString, "exponent overflows uint256") {
					t.Errorf("Expected overflow panic")
				}
			}()
		}

		actual := calculateTargetNBits(fixture.nBits)

		if expected != nil && actual.Cmp(expected) != 0 {
			t.Errorf("target NBits did not calculate correctly\nWanted 0x%x\nGot    0x%x", expected, actual)
			return
		}
	}
}
