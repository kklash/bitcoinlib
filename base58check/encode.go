// Package base58check provides APIs to encode and decode
// binary data in Bitcoin base58-check format.
package base58check

import (
	"bytes"
	"encoding/binary"

	"github.com/kklash/bitcoinlib/base58"
	"github.com/kklash/bitcoinlib/bhash"
)

// EncodeBase58CheckVersion encodes the given data into Base58-check, prepending
// a given version number onto it. If the version number is less than
// 0xff, it is encodes as a uint8, taking only one byte of space.
// If larger, it is encoded as a uint16, taking 2 bytes of space.
func EncodeVersion(data []byte, version uint16) string {
	buffer := new(bytes.Buffer)

	// Zcash and Decred version bytes are uint16s, others are singular bytes.
	var versionSized interface{} = version
	if version <= 0xff {
		versionSized = uint8(version)
	}
	binary.Write(buffer, binary.BigEndian, versionSized)
	buffer.Write(data)
	return Encode(buffer.Bytes())
}

// EncodeBase58Check encodes the given data as a base58-checksummed string.
func Encode(data []byte) string {
	hashed := bhash.DoubleSha256(data)
	return base58.Encode(append(data, hashed[:4]...))
}
