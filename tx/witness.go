package tx

import (
	"bytes"
	"errors"
	"io"

	"github.com/kklash/bitcoinlib/varint"
)

// Type Witness represents witness data for a single transaction Input.
type Witness [][]byte

const (
	// WitnessMaximumSize is an arbitrarily large limit for the size of a transaction
	// witness chunk, beyond which this library does not support decoding.
	WitnessMaximumSize = 0x20000000

	// WitnessChunkCountMaximum is an arbitrarily large limit for the number of witness chunks per
	// transaction, beyond which this library does not support decoding.
	WitnessChunkCountMaximum = 10000
)

var (
	// ErrInvalidWitnessFormat is returned when decoding a witness and the data is not valid.
	ErrInvalidWitnessFormat = errors.New("witness is not formatted correctly")

	// ErrIncompleteWitness is returned when trying to encode a witness where some chunks are nil.
	ErrIncompleteWitness = errors.New("cannot encode incomplete tx-witness")
)

// WitnessFromReader decodes a serialized Witness using data from reader.
// If the witness is not properly formatted, it returns ErrInvalidWitnessFormat.
func WitnessFromReader(reader io.Reader) (Witness, error) {
	witness, err := witnessFromReader(reader)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, ErrInvalidWitnessFormat
	} else if err != nil {
		return nil, err
	}

	return witness, nil
}

func witnessFromReader(reader io.Reader) (Witness, error) {
	nChunks, err := varint.FromReader(reader)
	if err != nil {
		return nil, err
	} else if nChunks > WitnessChunkCountMaximum {
		return nil, ErrInvalidWitnessFormat
	}

	witnessSize := uint64(0)
	witness := make(Witness, nChunks)

	for i := 0; varint.VarInt(i) < nChunks; i++ {
		chunkLength, err := varint.FromReader(reader)
		if err != nil {
			return nil, err
		}

		witnessSize += uint64(chunkLength)
		if witnessSize > WitnessMaximumSize {
			return nil, ErrInvalidWitnessFormat
		}

		witness[i] = make([]byte, chunkLength)
		_, err = io.ReadFull(reader, witness[i])
		if err != nil {
			return nil, err
		}
	}

	return witness, nil
}

// Size returns the serialized size of the Witness.
func (w Witness) Size() int {
	size := varint.VarInt(len(w)).Size()

	for _, chunk := range w {
		chunkLen := len(chunk)
		size += varint.VarInt(chunkLen).Size() + chunkLen
	}

	return size
}

// Bytes returns the serialized Witness as a byte-slice. Returns
// nil if any error occurs during serialization.
func (w Witness) Bytes() []byte {
	buf := new(bytes.Buffer)
	if _, err := w.WriteTo(buf); err != nil {
		return nil
	}

	return buf.Bytes()
}

// WriteTo implements the io.WriterTo interface. Writes the serialized Witness to a
// given io.Writer. Returns ErrIncompleteWitness if any of the witness chunks are nil.
func (w Witness) WriteTo(buf io.Writer) (n int64, err error) {
	for _, chunk := range w {
		if chunk == nil {
			err = ErrIncompleteWitness
			return
		}
	}

	nChunks := varint.VarInt(len(w))

	c, err := nChunks.WriteTo(buf)
	n += c
	if err != nil {
		return
	}

	for _, chunk := range w {
		chunkLen := varint.VarInt(len(chunk))

		c, err := chunkLen.WriteTo(buf)
		n += c
		if err != nil {
			return n, err
		}

		j, err := buf.Write(chunk)
		n += int64(j)
		if err != nil {
			return n, err
		}
	}

	return
}

// Clone returns a duplicate of the Witness.
func (w Witness) Clone() Witness {
	clone := make(Witness, len(w))
	for i, chunk := range w {
		if chunk != nil {
			clone[i] = make([]byte, len(chunk))
			copy(clone[i], chunk)
		}
	}

	return clone
}
