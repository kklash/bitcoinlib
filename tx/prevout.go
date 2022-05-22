package tx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/kklash/bitcoinlib/common"
)

// PrevOutSize is the serialized size of a PrevOut.
const PrevOutSize = 32 + 4

// Type PrevOut represents a previous transaction out-point. Used to encode Inputs.
type PrevOut struct {
	Hash  [32]byte
	Index uint32
}

func prevOutFromReader(reader io.Reader) (*PrevOut, error) {
	var (
		hash  [32]byte
		index uint32
	)

	_, err := io.ReadFull(reader, hash[:])
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &index)
	if err != nil {
		return nil, err
	}

	prevOut := &PrevOut{hash, index}
	return prevOut, nil
}

// Size returns the serialized size of the PrevOut. Always equal to PrevOutSize.
func (p *PrevOut) Size() int {
	return PrevOutSize
}

// Bytes returns the serialized PrevOut as a byte-slice. Returns
// nil if any error occurs during serialization.
func (p *PrevOut) Bytes() []byte {
	buf := new(bytes.Buffer)
	if _, err := p.WriteTo(buf); err != nil {
		return nil
	}

	return buf.Bytes()
}

// WriteTo implements the io.WriterTo interface. It writes the serialized PrevOut to a given io.Writer.
func (p *PrevOut) WriteTo(buf io.Writer) (n int64, err error) {
	j, err := buf.Write(p.Hash[:])
	n += int64(j)
	if err != nil {
		return
	}

	if err = binary.Write(buf, binary.LittleEndian, p.Index); err != nil {
		return
	}

	n += int64(binary.Size(p.Index))
	return
}

// Clone returns a pointer to a duplicate of the PrevOut.
func (p *PrevOut) Clone() *PrevOut {
	return &PrevOut{
		Hash:  p.Hash, // Copied by value
		Index: p.Index,
	}
}

// String returns the classic "txid:n" format reference of a PrevOut, where
// the txid is the reversed hash, and n is the previous output's index.
func (p *PrevOut) String() string {
	txid := common.ReverseBytes(p.Hash[:])
	return fmt.Sprintf("%x:%d", txid, p.Index)
}
