// Package tx provides a struct type which represents a Bitcoin transaction.
package tx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/common"
	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/bitcoinlib/varint"
)

const (
	// InputsMaximumCount is an approximation of the maximum possible number of inputs which
	// could be included in a transaction that would fit within a single block.
	InputsMaximumCount = constants.BlockMaxSize / InputMinimumSize

	// OutputsMaximumCount is an approximation of the maximum possible number of outputs which
	// could be included in a transaction that would fit within a single block.
	OutputsMaximumCount = constants.BlockMaxSize / OutputMinimumSize

	// MinimumSizeNoWitness is the minimum size of a transaction which Bitcoin Core will relay.
	//
	// See https://github.com/bitcoin/bitcoin/blob/5d32009f1a3b091299ff4a9345195b2359125f98/src/policy/policy.h#L25-L26
	MinimumSizeNoWitness = 82
)

// Tx represents a Bitcoin transaction.
type Tx struct {
	// Version represents a transaction version, which
	// can indicate compatibility or format changes.
	Version int32

	// Inputs are the sources of the Bitcoin being spent.
	Inputs []*Input

	// Outputs are the 'recipients' of the transaction.
	Outputs []*Output

	// Witnesses can be nil for non-segwit transactions.
	Witnesses []Witness

	// Locktime prevents the transaction from being mined until the current unix
	// time or block-height surpasses the locktime. If Locktime > 0x1dcd6500, it
	// is treated as a unix time stamp in seconds. Otherwise, it is treated as
	// a block height after which the transaction can be mined.
	Locktime uint32
}

var (
	// ErrIncompleteTx is returned when trying to encode a Tx when either of
	// tx.Inputs or tx.Outputs are nil. Also returned if any of the slice
	// contents of tx.Inputs, tx.Outputs, or tx.Witnesses are nil.
	ErrIncompleteTx = errors.New("cannot encode incomplete tx")

	// ErrInvalidTxFormat is returned when trying to decode a Tx which is not encoded properly.
	// Usually this is returned when the transaction is malformed and an
	// exceptionally large inputs/outputs count value is parsed.
	ErrInvalidTxFormat = errors.New("failed to decode improperly formatted transaction")
)

// FromReader decodes a serialized Bitcoin transaction from
// the given byte slice buf and returns a Tx.
func FromBytes(buf []byte) (*Tx, error) {
	return FromReader(bytes.NewReader(buf))
}

// FromReader decodes a serialized Bitcoin transaction
// from the given io.Reader r and returns a Tx.
func FromReader(reader io.Reader) (*Tx, error) {
	var version int32
	err := binary.Read(reader, binary.LittleEndian, &version)
	if err != nil {
		return nil, err
	}

	var witnessFlag [2]byte
	_, err = io.ReadFull(reader, witnessFlag[:])
	if err != nil {
		return nil, err
	}

	hasWitness := witnessFlag == constants.TxSegwitFlag
	if !hasWitness {
		reader = io.MultiReader(bytes.NewReader(witnessFlag[:]), reader)
	}

	nInputs, err := varint.FromReader(reader)
	if err != nil {
		return nil, err
	} else if nInputs > InputsMaximumCount {
		return nil, ErrInvalidTxFormat
	}

	inputs := make([]*Input, nInputs)
	for i := 0; varint.VarInt(i) < nInputs; i++ {
		inputs[i], err = InputFromReader(reader)
		if err != nil {
			return nil, err
		}
	}

	nOutputs, err := varint.FromReader(reader)
	if err != nil {
		return nil, err
	} else if nOutputs > OutputsMaximumCount {
		return nil, ErrInvalidTxFormat
	}

	outputs := make([]*Output, nOutputs)
	for i := 0; varint.VarInt(i) < nOutputs; i++ {
		outputs[i], err = OutputFromReader(reader)
		if err != nil {
			return nil, err
		}
	}

	var witnesses []Witness
	if hasWitness {
		witnesses = make([]Witness, nInputs)
		for i := 0; varint.VarInt(i) < nInputs; i++ {
			witnesses[i], err = WitnessFromReader(reader)
			if err != nil {
				return nil, err
			}
		}
	}

	var locktime uint32
	err = binary.Read(reader, binary.LittleEndian, &locktime)
	if err != nil {
		return nil, err
	}

	tx := &Tx{
		Version:   version,
		Inputs:    inputs,
		Outputs:   outputs,
		Witnesses: witnesses,
		Locktime:  locktime,
	}

	return tx, nil
}

// WriteTo implements the io.WriterTo interface. Writes the serialized Tx to
// the given io.Writer. If the Tx is not properly formatted, it returns ErrIncompleteTx.
func (tx *Tx) WriteTo(buf io.Writer) (int64, error) {
	n, err := tx.serialize(buf, true)
	if err != nil {
		return n, err
	}

	return n, nil
}

// WriteToNoWitness writes the serialized Tx to the given io.Writer, but
// does not include the witness flag or witness data with the transaction.
// If the Tx is not properly formatted, it returns ErrIncompleteTx.
func (tx *Tx) WriteToNoWitness(buf io.Writer) (int64, error) {
	n, err := tx.serialize(buf, false)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (tx *Tx) size(includeWitnesses bool) int {
	size := 4 + 4 // version + locktime

	if len(tx.Inputs) > 0 && includeWitnesses && len(tx.Witnesses) > 0 {
		size += len(constants.TxSegwitFlag)
	}

	size += varint.VarInt(len(tx.Inputs)).Size()
	for _, vin := range tx.Inputs {
		size += vin.Size()
	}

	size += varint.VarInt(len(tx.Outputs)).Size()
	for _, vout := range tx.Outputs {
		size += vout.Size()
	}

	if len(tx.Inputs) > 0 && includeWitnesses {
		for _, witness := range tx.Witnesses {
			size += witness.Size()
		}
	}

	return size
}

// Size returns the serialized size of the transaction, including witness data.
func (tx *Tx) Size() int {
	return tx.size(true)
}

// SizeNoWitness returns the serialized size of the transaction, excluding witness data.
func (tx *Tx) SizeNoWitness() int {
	return tx.size(false)
}

// Weight returns the number of block space weight units consumed by the transaction
// as defined by BIP141. Each byte of witness data counts as 1 weight unit. Each
// byte of non-witness data counts as 4 weight units.
func (tx *Tx) WeightUnits() int {
	legacyBytes := tx.SizeNoWitness()
	segwitBytes := tx.Size() - legacyBytes
	return legacyBytes*4 + segwitBytes
}

// VSize returns the number of virtual bytes consumed by the transaction
// for the purposes of fee estimation. A virtual byte (vbyte) is 4 weight units.
func (tx *Tx) VSize() int {
	vSize := float64(tx.WeightUnits()) / 4
	return int(math.Ceil(vSize))
}

// Bytes returns the serialized Tx as a byte-slice. Returns
// nil if any error occurs during serialization. Returns nil
// if an error was encountered during serialization.
func (tx *Tx) Bytes() []byte {
	return tx.bytes(true)
}

// Hex serializes the transaction and encodes it as a hex string.
// Returns "" if serialization fails for any reason.
func (tx *Tx) Hex() string {
	return hex.EncodeToString(tx.Bytes())
}

// String returns a debug string showing information about the transaction.
func (tx *Tx) String() string {
	return "<tx.Tx " +
		fmt.Sprintf("version=%d ", tx.Version) +
		fmt.Sprintf("n_inputs=%d ", len(tx.Inputs)) +
		fmt.Sprintf("n_outputs=%d ", len(tx.Outputs)) +
		fmt.Sprintf("n_witnesses=%d ", len(tx.Witnesses)) +
		fmt.Sprintf("locktime=%d ", tx.Locktime) +
		">"
}

// Bytes returns the serialized Tx as a byte-slice, but does not
// include the witness flag or witness data with the transaction. Returns nil
// if an error was encountered during serialization.
func (tx *Tx) BytesNoWitness() []byte {
	return tx.bytes(false)
}

func (tx *Tx) bytes(includeWitnesses bool) []byte {
	buf := new(bytes.Buffer)
	if _, err := tx.serialize(buf, includeWitnesses); err != nil {
		return nil
	}

	return buf.Bytes()
}

func (tx *Tx) serialize(buf io.Writer, includeWitnesses bool) (n int64, err error) {
	if !tx.canSerialize() {
		err = ErrIncompleteTx
		return
	}

	// Add version number
	if err = binary.Write(buf, binary.LittleEndian, tx.Version); err != nil {
		return
	}
	n += int64(binary.Size(tx.Version))

	// Add witness flag if needed
	if len(tx.Inputs) > 0 && includeWitnesses && len(tx.Witnesses) > 0 {
		j, err := buf.Write(constants.TxSegwitFlag[:])
		n += int64(j)
		if err != nil {
			return n, err
		}
	}

	// Serialize inputs
	nInputs := varint.VarInt(len(tx.Inputs))

	c, err := nInputs.WriteTo(buf)
	n += c
	if err != nil {
		return
	}

	for _, vin := range tx.Inputs {
		c, err = vin.WriteTo(buf)
		n += c
		if err != nil {
			return
		}
	}

	// Serialize outputs
	nOutputs := varint.VarInt(len(tx.Outputs))

	c, err = nOutputs.WriteTo(buf)
	n += c
	if err != nil {
		return
	}

	for _, vout := range tx.Outputs {
		c, err = vout.WriteTo(buf)
		n += c
		if err != nil {
			return
		}
	}

	// Serialize witnesses
	if len(tx.Inputs) > 0 && includeWitnesses {
		for _, witness := range tx.Witnesses {
			c, err = witness.WriteTo(buf)
			n += c
			if err != nil {
				return
			}
		}
	}

	if err = binary.Write(buf, binary.LittleEndian, tx.Locktime); err != nil {
		return
	}

	n += int64(binary.Size(tx.Locktime))
	return
}

func (tx *Tx) canSerialize() bool {
	if tx.Inputs == nil || tx.Outputs == nil {
		return false
	}

	for _, vin := range tx.Inputs {
		if vin == nil {
			return false
		}
	}

	for _, vout := range tx.Outputs {
		if vout == nil {
			return false
		}
	}

	if tx.Witnesses != nil {
		if len(tx.Witnesses) != len(tx.Inputs) {
			return false
		}

		for _, witness := range tx.Witnesses {
			if witness == nil {
				return false
			}
		}
	}

	return true
}

// Hash returns the double SHA256 hash of the serialized transaction,
// optionally including witness data for a witness hash.
func (tx *Tx) Hash(includeWitness bool) ([32]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := tx.serialize(buf, includeWitness); err != nil {
		return [32]byte{}, err
	}

	return bhash.DoubleSha256(buf.Bytes()), nil
}

// Id returns the TXID of the transaction, optionally
// including witness data for a witness TXID.
func (tx *Tx) Id(includeWitness bool) (string, error) {
	doubleHashArr, err := tx.Hash(includeWitness)
	if err != nil {
		return "", err
	}

	// Reverse the hash for a TXID
	txidBytes := common.ReverseBytes(doubleHashArr[:])
	return hex.EncodeToString(txidBytes), nil
}
