package varint

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func hex2bytes(h string) []byte {
	buf, _ := hex.DecodeString(h)
	return buf
}

func ExampleVarInt_encoding() {
	vi, err := FromNumber(0xabd280)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%x\n", vi.Bytes())
	// Output: fe80d2ab00
}

func TestVarInt(t *testing.T) {
	byteSlices := [][]byte{
		hex2bytes("fdfe00"),
		hex2bytes("2d"),
		hex2bytes("a3"),
		hex2bytes("fd2301"),
		hex2bytes("ff758277472d8ac72b"),
		hex2bytes("ff1388559977000000"),
		hex2bytes("fe9374bd33"),
		hex2bytes("fe0010cc82"),
	}

	uintValues := []VarInt{
		0xfe,
		0x2d,
		0xa3,
		0x123,
		0x2bc78a2d47778275,
		0x7799558813,
		0x33bd7493,
		0x82cc1000,
	}

	for i := 0; i < len(byteSlices); i++ {
		output, err := FromReader(bytes.NewBuffer(byteSlices[i]))
		if err != nil {
			t.Errorf(err.Error())
		} else if output != uintValues[i] {
			t.Errorf("VarInt did not decode as expected - Wanted %d, got %d", uintValues[i], output)
		}

		encoded := uintValues[i].Bytes()
		if !bytes.Equal(encoded, byteSlices[i]) {
			t.Errorf("VarInt did not encode as expected - Wanted %x, got %x", byteSlices[i], encoded)
		} else if size := uintValues[i].Size(); size != len(encoded) {
			t.Errorf("VarInt size does not match - Wanted %d, got %d", len(encoded), size)
		}
	}
}
