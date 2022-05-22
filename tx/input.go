package tx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/bitcoinlib/varint"
)

const (
	// InputMinimumSize is the smallest possible size for a transaction input.
	InputMinimumSize = PrevOutSize + 1 + 4
)

var (
	// ErrInvalidInputFormat is returned when decoding an input and the data is not valid.
	ErrInvalidInputFormat = errors.New("tx-input is not formatted correctly")

	// ErrIncompleteInput is returned when trying to encode an input with no script or PrevOut.
	ErrIncompleteInput = errors.New("cannot encode incomplete tx-input")
)

// Input represents a transaction input.
type Input struct {
	PrevOut  *PrevOut
	Script   []byte
	Sequence uint32
}

// InputFromReader decodes a serialized Input using data from reader.
// If the input is not properly formatted, it returns ErrInvalidInputFormat.
func InputFromReader(reader io.Reader) (*Input, error) {
	input, err := inputFromReader(reader)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, ErrInvalidInputFormat
	} else if err != nil {
		return nil, err
	}

	return input, nil
}

func inputFromReader(reader io.Reader) (*Input, error) {
	prevOut, err := prevOutFromReader(reader)
	if err != nil {
		return nil, err
	}

	scriptLen, err := varint.FromReader(reader)
	if err != nil {
		return nil, err
	} else if scriptLen > constants.BlockMaxSize {
		return nil, ErrInvalidInputFormat
	}

	script := make([]byte, scriptLen)
	_, err = io.ReadFull(reader, script)
	if err != nil {
		return nil, err
	}

	var sequence uint32
	err = binary.Read(reader, binary.LittleEndian, &sequence)
	if err != nil {
		return nil, err
	}

	input := &Input{
		PrevOut:  prevOut,
		Script:   script,
		Sequence: sequence,
	}

	return input, nil
}

// Size returns the serialized size of the input.
func (i *Input) Size() int {
	scriptLen := len(i.Script)
	return i.PrevOut.Size() + 4 + scriptLen + varint.VarInt(scriptLen).Size()
}

// Bytes returns the serialized Input as a byte-slice. Returns
// nil if any error occurs during serialization.
func (i *Input) Bytes() []byte {
	buf := new(bytes.Buffer)
	if _, err := i.WriteTo(buf); err != nil {
		return nil
	}

	return buf.Bytes()
}

// WriteTo implements the io.WriterTo interface. Writes the serialized Input to a
// given io.Writer. Returns ErrIncompleteInput if either i.Script or i.PrevOut are nil.
func (i *Input) WriteTo(buf io.Writer) (n int64, err error) {
	if i.PrevOut == nil || i.Script == nil {
		err = ErrIncompleteInput
		return
	}

	c, err := i.PrevOut.WriteTo(buf)
	n += c
	if err != nil {
		return
	}

	scriptLen := varint.VarInt(len(i.Script))

	c, err = scriptLen.WriteTo(buf)
	n += c
	if err != nil {
		return
	}

	j, err := buf.Write(i.Script)
	n += int64(j)
	if err != nil {
		return
	}

	if err = binary.Write(buf, binary.LittleEndian, i.Sequence); err != nil {
		return
	}

	n += int64(binary.Size(i.Sequence))
	return
}

// Clone returns a pointer to a duplicate of the Input.
func (i *Input) Clone() *Input {
	clone := &Input{
		PrevOut:  i.PrevOut.Clone(),
		Sequence: i.Sequence,
	}

	if i.Script != nil {
		clone.Script = make([]byte, len(i.Script))
		copy(clone.Script, i.Script)
	}

	return clone
}
