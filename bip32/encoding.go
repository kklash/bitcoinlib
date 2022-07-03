package bip32

import (
	"encoding/binary"
)

func serialize32(v uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}
