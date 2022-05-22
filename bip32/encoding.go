package bip32

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

func serialize32(v uint32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, v)
	return buf.Bytes()
}

func parse32(b []byte) uint32 {
	buf := bytes.NewBuffer(b)
	var v uint32
	if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
		panic("failed to parse uint32 from bytes: " + err.Error())
	}
	return v
}

func parse256(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}
