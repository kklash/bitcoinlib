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
	// OutputMinimumSize is the smallest possible size for a transaction output.
	OutputMinimumSize = 8 + 1
)

var (
	// ErrInvalidOutputFormat is returned when decoding an output and the data is not valid.
	ErrInvalidOutputFormat = errors.New("tx-output is not formatted correctly")

	// ErrIncompleteOutput is returned when trying to encode an output with no script.
	ErrIncompleteOutput = errors.New("cannot encode incomplete tx-output")
)

// Output represents a transaction output.
type Output struct {
	Value  uint64
	Script []byte
}

// OutputFromReader decodes a serialized Output using data from reader.
// If the output is not properly formatted, it returns ErrInvalidOutputFormat.
func OutputFromReader(reader io.Reader) (*Output, error) {
	output, err := outputFromReader(reader)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, ErrInvalidOutputFormat
	} else if err != nil {
		return nil, err
	}

	return output, nil
}

func outputFromReader(reader io.Reader) (*Output, error) {
	var value uint64
	err := binary.Read(reader, binary.LittleEndian, &value)
	if err != nil {
		return nil, err
	}

	scriptLength, err := varint.FromReader(reader)
	if err != nil {
		return nil, err
	} else if scriptLength > constants.BlockMaxSize {
		return nil, ErrInvalidOutputFormat
	}

	script := make([]byte, scriptLength)
	_, err = io.ReadFull(reader, script)
	if err != nil {
		return nil, err
	}

	output := &Output{
		Value:  value,
		Script: script,
	}
	return output, nil
}

// Size returns the serialized size of the output.
func (o *Output) Size() int {
	scriptLen := len(o.Script)
	return 8 + varint.VarInt(scriptLen).Size() + scriptLen
}

// Bytes returns the serialized Output as a byte-slice. Returns
// nil if any error occurs during serialization.
func (o *Output) Bytes() []byte {
	buf := new(bytes.Buffer)
	if _, err := o.WriteTo(buf); err != nil {
		return nil
	}

	return buf.Bytes()
}

// WriteTo implements the io.WriterTo interface. Writes the serialized
// Output to a given io.Writer. Returns ErrIncompleteOutput if o.Script is nil.
func (o *Output) WriteTo(buf io.Writer) (n int64, err error) {
	if o.Script == nil {
		err = ErrIncompleteOutput
		return
	}

	if err = binary.Write(buf, binary.LittleEndian, o.Value); err != nil {
		return
	}
	n += int64(binary.Size(o.Value))

	scriptLen := varint.VarInt(len(o.Script))

	c, err := scriptLen.WriteTo(buf)
	n += c
	if err != nil {
		return
	}

	j, err := buf.Write(o.Script)
	n += int64(j)
	if err != nil {
		return
	}

	return
}

// Clone returns a pointer to a duplicate of the Output.
func (o *Output) Clone() *Output {
	clone := &Output{Value: o.Value}

	if o.Script != nil {
		clone.Script = make([]byte, len(o.Script))
		copy(clone.Script, o.Script)
	}

	return clone
}
