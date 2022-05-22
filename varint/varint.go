// Package varint exposes a type VarInt which can be used for convenient
// encoding and decoding of variable-sized integers in the Bitcoin protocol.
package varint

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

// VarInt represents an unsigned integer which can be serialized according to
// the standards of the bitcoin protocol.
type VarInt uint64

// FromReader returns a new VarInt read from the given io.Reader.
// You can determine how many bytes were read by calling Size() on the returned VarInt.
func FromReader(reader io.Reader) (VarInt, error) {
	firstBytes := make([]byte, 1)
	_, err := io.ReadFull(reader, firstBytes)
	if err != nil {
		return 0, err
	}

	var vi VarInt
	switch firstBytes[0] {
	case 0xff:
		var size uint64
		err = binary.Read(reader, binary.LittleEndian, &size)
		vi = VarInt(size)
	case 0xfe:
		var size uint32
		err = binary.Read(reader, binary.LittleEndian, &size)
		vi = VarInt(size)
	case 0xfd:
		var size uint16
		err = binary.Read(reader, binary.LittleEndian, &size)
		vi = VarInt(size)
	default:
		return VarInt(firstBytes[0]), nil
	}
	if err != nil {
		return 0, err
	}

	return vi, nil
}

// FromBytes attempts to read a VarInt from the given byte slice. It returns
// a VarInt, and any error encountered. This function's use is not recommended
// because it does not return the number of bytes read. Use FromReader instead
// so that the caller can track their index when reading a stream of data.
func FromBytes(buf []byte) (VarInt, error) {
	return FromReader(bytes.NewReader(buf))
}

// FromNumber returns a new VarInt constructed from the given unsigned integer n.
// n must be a uint, uint8, uint16, uint32, or uint64, OR if it is a signed integer
// it must be greater than 0.
func FromNumber(n interface{}) (VarInt, error) {
	switch n.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return VarInt(reflect.ValueOf(n).Uint()), nil
	case int, int8, int16, int32, int64:
		value := reflect.ValueOf(n).Int()
		if value < 0 {
			return 0, fmt.Errorf("cannot create varint from value less than zero: %d", value)
		}
		return VarInt(value), nil
	default:
		return 0, fmt.Errorf("%v is not a number, cannot make varint", reflect.TypeOf(n))
	}
}

// Size returns the encoded length of a VarInt in bytes.
func (v VarInt) Size() int {
	switch {
	case v > 0xffffffff:
		return 9
	case v > 0xffff:
		return 5
	case v > 0xfc:
		return 3
	default:
		return 1
	}
}

// WriteTo implements the io.WriterTo interface, encoding the VarInt
// and writing it to the given writer w.
func (v VarInt) WriteTo(w io.Writer) (n int64, err error) {
	var (
		prefixByte byte
		sized      interface{}
	)

	switch {
	case v > 0xffffffff:
		prefixByte = 0xff
		sized = uint64(v)
	case v > 0xffff:
		prefixByte = 0xfe
		sized = uint32(v)
	case v > 0xfc:
		prefixByte = 0xfd
		sized = uint16(v)
	default:
		prefixByte = byte(v)
		sized = nil
	}

	if _, err = w.Write([]byte{prefixByte}); err != nil {
		return
	}
	n++

	if sized != nil {
		if err = binary.Write(w, binary.LittleEndian, sized); err != nil {
			return
		}
		n += int64(binary.Size(sized))
	}

	return
}

// Bytes returns the serialized VarInt as a byte-slice. Returns
// nil if any error occurs during serialization.
func (v VarInt) Bytes() []byte {
	buf := new(bytes.Buffer)
	if _, err := v.WriteTo(buf); err != nil {
		return nil
	}

	return buf.Bytes()
}
