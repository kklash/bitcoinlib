// Package script exposes functions for creating Bitcoin output scripts.
package script

import (
	"bytes"
	"errors"
	"io"

	"github.com/kklash/bitcoinlib/constants"
)

var (
	// ErrInvalidScript is returned by script decoding functions if the provided
	// script does not match the expected script pub key format.
	ErrInvalidScript = errors.New("cannot decode improperly formatted script")

	// ErrInvalidPublicKeyLength indicates an improper length of public key was passed
	// to an address-making function.
	ErrInvalidPublicKeyLength = errors.New("invalid public key length")

	// ErrNotPushOnlyScript is returned if a function expects a push-only script,
	// and receives one which includes other non-push opcodes.
	ErrNotPushOnlyScript = errors.New("expected script which is push-only")
)

// ClassifyOutput determines the type of the given script pub key, returning an address format string.
// If the script type is not recognized, we return constants.FormatNONSTANDARD.
func ClassifyOutput(script []byte) constants.AddressFormat {
	switch {
	case IsP2PKH(script):
		return constants.FormatP2PKH
	case IsP2SH(script):
		return constants.FormatP2SH
	case IsP2WPKH(script):
		return constants.FormatP2WPKH
	case IsP2WSH(script):
		return constants.FormatP2WSH
	default:
		return constants.FormatNONSTANDARD
	}
}

// parsePushIntOpCode parses the given opcode as an integer-pushing operation.
// It returns the pushed integer as a byte, which could be 0 - 16 or 0x81
// (which represents negative 1). It also returns an 'ok' boolean which is
// false if the op code is not an integer-pushing operation.
func parsePushIntOpCode(op byte) (byte, bool) {
	var b byte
	switch {
	case op == constants.OP_0:
		b = 0
	case op == constants.OP_1NEGATE:
		b = 0x81 // AKA -1
	case op >= constants.OP_1 && op <= constants.OP_16:
		b = op - constants.OP_1 + 1
	default:
		return 0, false
	}

	return b, true
}

// Decompile breaks down a script into chunks of op codes and pushdata blocks.
// The return type is []Union<byte, []byte>.
func Decompile(script []byte) ([]interface{}, error) {
	chunks := make([]interface{}, 0, 8)

	r := bytes.NewReader(script)

	for {
		nextByte, err := r.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if nextByte > 0 && nextByte <= constants.OP_PUSHDATA4 {
			r.Seek(-1, io.SeekCurrent) // backtrack to include nextByte
			buf, err := ReadData(r)
			if err != nil {
				return nil, err
			}

			chunks = append(chunks, buf)
		} else {
			chunks = append(chunks, nextByte)
		}
	}

	return chunks, nil
}

// Stackify turns a push-only script into a slice of byte-slices, with one
// slice for each piece of data or number that was pushed to the stack by the script.
// Returns ErrNotPushOnlyScript if the given script is not a push-only script.
// This function is used to turn unlocking scripts into segregated witness script signatures.
func Stackify(script []byte) ([][]byte, error) {
	chunks, err := Decompile(script)
	if err != nil {
		return nil, err
	}

	stack := make([][]byte, len(chunks))

	for i, chunk := range chunks {
		switch chunk.(type) {
		case byte:
			pushedByte, ok := parsePushIntOpCode(chunk.(byte))
			if !ok {
				return nil, ErrNotPushOnlyScript
			}

			// An empty byte vector is treated as zero (false) for stack operations
			if pushedByte == 0 {
				stack[i] = []byte{}
			} else {
				stack[i] = []byte{pushedByte}
			}

		case []byte:
			stack[i] = chunk.([]byte)
		}
	}

	return stack, nil
}

// TODO need more test vectors for this
func StripOpCode(script []byte, op byte) ([]byte, error) {
	chunks, err := Decompile(script)
	if err != nil {
		return nil, err
	}

	recompiled := new(bytes.Buffer)
	for _, chunk := range chunks {
		switch chunk.(type) {
		case byte:
			b := chunk.(byte)
			if b != op {
				recompiled.WriteByte(b)
			}
		case []byte:
			push := PushData(chunk.([]byte))
			recompiled.Write(push)
		}
	}

	if recompiled.Len() == 0 {
		return []byte{}, nil
	}

	return recompiled.Bytes(), nil
}
