// Package blocks provides a struct type which represents a Bitcoin block.
package blocks

import (
	"bytes"
	"errors"
	"io"

	"github.com/kklash/bitcoinlib/blocks/blockheader"
	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/bitcoinlib/tx"
	"github.com/kklash/bitcoinlib/varint"
)

var (
	// ErrInvalidFormat is returned when decoding an improperly serialized block.
	ErrInvalidFormat = errors.New("block is not formatted correctly")
)

// Block represents a bitcoin block, including the header.
type Block struct {
	Header       *blockheader.BlockHeader
	Transactions []*tx.Tx
}

// FromReader decodes a serialized BlockHeader using data from reader.
// If the block is not properly formatted, it returns ErrInvalidFormat.
func FromReader(reader io.Reader) (*Block, error) {
	block, err := fromReader(reader)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, ErrInvalidFormat
	} else if err != nil {
		return nil, err
	}

	return block, nil
}

func fromReader(r io.Reader) (*Block, error) {
	header, err := blockheader.FromReader(r)
	if err != nil {
		return nil, err
	}

	nTx, err := varint.FromReader(r)
	if err != nil {
		return nil, err
	} else if int(nTx) > constants.BlockMaxSize/tx.MinimumSizeNoWitness {
		return nil, ErrInvalidFormat
	}

	transactions := make([]*tx.Tx, nTx)
	for i := 0; i < int(nTx); i++ {
		transactions[i], err = tx.FromReader(r)
		if err != nil {
			return nil, err
		}
	}

	block := &Block{header, transactions}
	return block, nil
}

// WriteTo implements the io.WriterTo interface. Writes the serialized Block to
// the given io.Writer.
func (block *Block) WriteTo(w io.Writer) (n int64, err error) {
	c, err := block.Header.WriteTo(w)
	n += c
	if err != nil {
		return n, err
	}

	nTx := varint.VarInt(len(block.Transactions))
	c, err = nTx.WriteTo(w)
	n += c
	if err != nil {
		return n, err
	}

	for _, txn := range block.Transactions {
		c, err := txn.WriteTo(w)
		n += c
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// Size returns the serialized size of the block including witness data.
func (block *Block) Size() int {
	size := block.Header.Size()
	size += varint.VarInt(len(block.Transactions)).Size()

	for _, txn := range block.Transactions {
		size += txn.Size()
	}

	return size
}

// WeightUnits returns the total weight of the block. Should never
// be above 4,000,000 (4 mega-weight-units) as per BIP141
func (block *Block) WeightUnits() int {
	weight := block.Header.Size() * 4
	weight += varint.VarInt(len(block.Transactions)).Size() * 4

	for _, txn := range block.Transactions {
		weight += txn.WeightUnits()
	}

	return weight
}

// Bytes returns the serialized Block as a byte-slice. Returns
// nil if any error occurs during serialization.
func (block *Block) Bytes() []byte {
	buf := new(bytes.Buffer)
	if _, err := block.WriteTo(buf); err != nil {
		return nil
	}

	return buf.Bytes()
}
