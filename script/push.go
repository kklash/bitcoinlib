package script

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/kklash/bitcoinlib/constants"
)

var (
	// ErrGreaterThanMaxInt is returned by ReadNumber if the script being
	// decoded has encoded a number which overflows int64.
	ErrGreaterThanMaxInt = errors.New("cannot decode number larger than maximum int64")
)

// PushData returns a byte slices with the necessary opcodes prepended to
// push the given data to the stack. Note that you cannot push more than
// 520 bytes to the stack, but since this limit does not apply to
// witnesses, this function doesn't enforce that limit.
func PushData(data []byte) []byte {
	dataSize := len(data)
	if dataSize <= constants.OP_DATA_75 { // 0x00 - 0x4b are direct push opcodes
		return append([]byte{byte(dataSize)}, data...)
	}

	script := new(bytes.Buffer)
	switch {
	case dataSize <= 0xff:
		script.WriteByte(constants.OP_PUSHDATA1)
		binary.Write(script, binary.LittleEndian, uint8(dataSize))
	case dataSize <= 0xffff:
		script.WriteByte(constants.OP_PUSHDATA2)
		binary.Write(script, binary.LittleEndian, uint16(dataSize))
	case dataSize <= 0xffffffff:
		script.WriteByte(constants.OP_PUSHDATA4)
		binary.Write(script, binary.LittleEndian, uint32(dataSize))
	default:
		panic("cannot push data larger than max uint32")
	}

	script.Write(data)
	return script.Bytes()
}

// ReadData attempts to read a byte array from the given reader r. It expects
// first byte from the reader to either be a push opcode, or if the array
// is less than 0x4c bytes long, the byte array length itself.
func ReadData(r io.Reader) ([]byte, error) {
	firstBytes := make([]byte, 1)
	_, err := io.ReadFull(r, firstBytes)
	if err != nil {
		return nil, err
	}

	var dataSize uint32

	if firstBytes[0] <= constants.OP_DATA_75 {
		dataSize = uint32(firstBytes[0])
	} else if firstBytes[0] == constants.OP_PUSHDATA1 {
		var size uint8
		err = binary.Read(r, binary.LittleEndian, &size)
		dataSize = uint32(size)
	} else if firstBytes[0] == constants.OP_PUSHDATA2 {
		var size uint16
		err = binary.Read(r, binary.LittleEndian, &size)
		dataSize = uint32(size)
	} else if firstBytes[0] == constants.OP_PUSHDATA4 {
		err = binary.Read(r, binary.LittleEndian, &dataSize)
	} else {
		return nil, fmt.Errorf("received unexpected byte '%d' when trying to read data", firstBytes[0])
	}

	if err != nil {
		return nil, err
	}

	buf := make([]byte, dataSize)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// PushNumber returns a PushData script which
// pushes a given number onto the stack.
func PushNumber(n int64) []byte {
	switch {
	case n == -1:
		return []byte{constants.OP_1NEGATE}
	case n == 0:
		return []byte{0}
	case n >= 1 && n <= 16:
		return []byte{byte(n + constants.OP_1 - 1)}
	}

	/***** This section copied from btcsuite/btcd/txscript/scriptnum.go *****/
	// Take the absolute value and keep track of whether it was originally
	// negative.
	isNegative := n < 0
	if isNegative {
		n = -n
	}

	// Encode to little endian.
	result := make([]byte, 0, 8)
	for n > 0 {
		result = append(result, byte(n&0xff))
		n >>= 8
	}

	// When the most significant byte already has the high bit set, an
	// additional high byte is required to indicate whether the number is
	// negative or positive.  The additional byte is removed when converting
	// back to an integral and its high bit is used to denote the sign.
	//
	// Otherwise, when the most significant byte does not already have the
	// high bit set, use it to indicate the value is negative, if needed.
	if result[len(result)-1]&0x80 != 0 {
		if isNegative {
			result = append(result, 0x80)
		} else {
			result = append(result, 0x00)
		}
	} else if isNegative {
		result[len(result)-1] |= 0x80
	}
	/***** This section copied from btcsuite/btcd/txscript/scriptnum.go *****/

	return PushData(result)
}

// ReadNumber attemps to read an int64 from the script reader r and
// returns it, along with any error encountered. Returns
// ErrGreaterThanMaxInt if the data read from the script overflows int64.
func ReadNumber(r io.Reader) (int64, error) {
	firstBytes := make([]byte, 1)
	_, err := io.ReadFull(r, firstBytes)
	if err != nil {
		return 0, err
	}

	if firstBytes[0] == constants.OP_1NEGATE {
		return -1, nil
	} else if firstBytes[0] == 0 {
		return 0, nil
	} else if int(firstBytes[0])-0x50 > 0 && int(firstBytes[0])-0x50 < 17 {
		return int64(firstBytes[0] - 0x50), nil
	}

	// re-read firstBytes
	r = io.MultiReader(bytes.NewReader(firstBytes), r)

	numBytes, err := ReadData(r)
	if err != nil {
		return 0, err
	}

	if len(numBytes) > 8 {
		return 0, ErrGreaterThanMaxInt
	}

	var (
		isNegative bool
		value      int64
	)

	for i := 0; i < len(numBytes); i++ {
		byteValue := numBytes[i]
		if i == len(numBytes)-1 && byteValue&0x80 != 0 {
			byteValue -= 0x80
			isNegative = true
		}
		value += int64(byteValue) << (8 * uint(i))
	}

	if isNegative {
		value = -value
	}

	return value, nil
}
