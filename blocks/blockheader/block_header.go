// Package blockheader provides a struct type which represents a Bitcoin block header.
package blockheader

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math/big"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/common"
)

// BlockHeaderSize is the byte size of a serialized block header.
const BlockHeaderSize = (4 * 4) + (32 * 2)

// BlockHeader is a struct representing the header portion of a block.
type BlockHeader struct {
	Version            int32
	PreviousHeaderHash [32]byte
	MerkleRootHash     [32]byte
	Time               uint32
	NBits              uint32
	Nonce              uint32
}

var (
	// ErrInvalidFormat is returned when decoding a block header and the data is not valid.
	ErrInvalidFormat = errors.New("block header is not formatted correctly")
)

// FromReader decodes a serialized BlockHeader using data from reader.
// If the header is not properly formatted, it returns ErrInvalidFormat.
func FromReader(reader io.Reader) (*BlockHeader, error) {
	header, err := fromReader(reader)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, ErrInvalidFormat
	} else if err != nil {
		return nil, err
	}

	return header, nil
}

func fromReader(reader io.Reader) (*BlockHeader, error) {
	header := new(BlockHeader)

	if err := binary.Read(reader, binary.LittleEndian, &header.Version); err != nil {
		return nil, err
	}

	if _, err := io.ReadFull(reader, header.PreviousHeaderHash[:]); err != nil {
		return nil, err
	}
	common.ReverseBytesInPlace(header.PreviousHeaderHash[:])

	if _, err := io.ReadFull(reader, header.MerkleRootHash[:]); err != nil {
		return nil, err
	}
	common.ReverseBytesInPlace(header.MerkleRootHash[:])

	if err := binary.Read(reader, binary.LittleEndian, &header.Time); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &header.NBits); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &header.Nonce); err != nil {
		return nil, err
	}

	return header, nil
}

// WriteTo implements the io.WriterTo interface. Writes the serialized BlockHeader
// to the given io.Writer.
func (header *BlockHeader) WriteTo(buf io.Writer) (n int64, err error) {
	if err = binary.Write(buf, binary.LittleEndian, header.Version); err != nil {
		return
	}
	n += int64(binary.Size(header.Version))

	c, err := buf.Write(common.ReverseBytes(header.PreviousHeaderHash[:]))
	if err != nil {
		return
	}
	n += int64(c)

	c, err = buf.Write(common.ReverseBytes(header.MerkleRootHash[:]))
	if err != nil {
		return
	}
	n += int64(c)

	if err = binary.Write(buf, binary.LittleEndian, header.Time); err != nil {
		return
	}
	n += int64(binary.Size(header.Time))

	if err = binary.Write(buf, binary.LittleEndian, header.NBits); err != nil {
		return
	}
	n += int64(binary.Size(header.NBits))

	if err = binary.Write(buf, binary.LittleEndian, header.Nonce); err != nil {
		return
	}
	n += int64(binary.Size(header.Nonce))

	return
}

// Size returns the serialized size of the BlockHeader. Always equal to BlockHeaderSize.
func (header *BlockHeader) Size() int {
	return BlockHeaderSize
}

// Bytes returns the serialized BlockHeader as a byte-slice. Returns
// nil if any error occurs during serialization.
func (header *BlockHeader) Bytes() []byte {
	buf := new(bytes.Buffer)
	if _, err := header.WriteTo(buf); err != nil {
		return nil
	}

	return buf.Bytes()
}

// TargetNBits returns the 256-bit nBits target value after decoding header.NBits from compact notation.
func (header *BlockHeader) TargetNBits() *big.Int {
	return calculateTargetNBits(header.NBits)
}

// Hash returns the reversed double SHA256 hash of the serialized block header,
// and any error encountered during serialization.
func (header *BlockHeader) Hash() ([32]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := header.WriteTo(buf); err != nil {
		return [32]byte{}, err
	}

	doubleHashed := bhash.DoubleSha256(buf.Bytes())
	common.ReverseBytesInPlace(doubleHashed[:])
	return doubleHashed, nil
}
