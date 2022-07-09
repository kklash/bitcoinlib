// Package unspent provides structures for managing groups of Bitcoin unspent transaction outputs (UTXOs).
package unspent

import (
	"bytes"
	"encoding/hex"

	"github.com/kklash/bitcoinlib/blocks"
	"github.com/kklash/bitcoinlib/common"
	"github.com/kklash/bitcoinlib/tx"
)

// OutputSet represents a set of unspent transaction outputs,
// mapped by their transaction outpoints for faster referencing.
type OutputSet struct {
	byOutpoint map[tx.PrevOut]*Output
}

// NewOutputSet returns a new OutputSet from an array of unspent Output structs.
func NewOutputSet(outputs []*Output) *OutputSet {
	unspentOutputs := new(OutputSet)
	for _, output := range outputs {
		unspentOutputs.AddOutput(output)
	}
	return unspentOutputs
}

// Size returns the number of unspent outputs in the OutputSet.
func (unspentOutputs *OutputSet) Size() int {
	return len(unspentOutputs.byOutpoint)
}

// Clone returns a pointer to a duplicate of the OutputSet.
func (unspentOutputs *OutputSet) Clone() *OutputSet {
	clone := new(OutputSet)
	for _, output := range unspentOutputs.byOutpoint {
		clone.AddOutput(output)
	}
	return clone
}

// Slice allocates and returns the OutputSet in the form of a slice of Output pointers.
func (unspentOutputs *OutputSet) Slice() []*Output {
	i := 0
	utxos := make([]*Output, unspentOutputs.Size())

	for _, output := range unspentOutputs.byOutpoint {
		utxos[i] = output
		i += 1
	}

	return utxos
}

// GetByOutpoint looks up an unspent output by its transaction outpoint.
// Returns nil if not found.
func (unspentOutputs *OutputSet) GetByOutpoint(outpoint *tx.PrevOut) *Output {
	if outpoint == nil || unspentOutputs.Size() == 0 {
		return nil
	}

	output, ok := unspentOutputs.byOutpoint[*outpoint]
	if ok {
		return output
	}

	return nil
}

// GetByHash looks up an unspent output by its transaction outpoint
// hash and index. Returns nil if not found.
func (unspentOutputs *OutputSet) GetByHash(hash [32]byte, index uint32) *Output {
	outpoint := &tx.PrevOut{Hash: hash, Index: index}
	return unspentOutputs.GetByOutpoint(outpoint)
}

// GetByTxid looks up an unspent output by its transaction outpoint
// TXID string and index. Returns nil if not found. The TXID string
// is the byte-order reversal of the TX hash.
func (unspentOutputs *OutputSet) GetByTxid(txid string, index uint32) *Output {
	if len(txid) != 64 || unspentOutputs.Size() == 0 {
		return nil
	}

	txidBytes, err := hex.DecodeString(txid)
	if err != nil {
		return nil
	}

	var hash [32]byte
	copy(hash[:], txidBytes)
	common.ReverseBytesInPlace(hash[:])

	outpoint := &tx.PrevOut{Hash: hash, Index: index}
	return unspentOutputs.GetByOutpoint(outpoint)
}

// AddOutput adds an unspent output to the OutputSet,
// overwriting any existing outputs, if present.
func (unspentOutputs *OutputSet) AddOutput(output *Output) {
	if unspentOutputs.byOutpoint == nil {
		unspentOutputs.byOutpoint = make(map[tx.PrevOut]*Output)
	}

	unspentOutputs.byOutpoint[*output.Outpoint] = output
}

// RemoveByOutpoint removes an unspent output from the OutputSet
// as specified by the given transaction outpoint.
func (unspentOutputs *OutputSet) RemoveByOutpoint(outpoint *tx.PrevOut) {
	if unspentOutputs.byOutpoint != nil {
		delete(unspentOutputs.byOutpoint, *outpoint)
	}
}

// RemoveByHash removes an unspent output from the OutputSet matching the
// given TX hash and output index.
func (unspentOutputs *OutputSet) RemoveByHash(hash [32]byte, index uint32) {
	if unspentOutputs.byOutpoint != nil {
		outpoint := &tx.PrevOut{Hash: hash, Index: index}
		unspentOutputs.RemoveByOutpoint(outpoint)
	}
}

// RemoveByHash removes an unspent output from the OutputSet matching the
// given TXID and output index.
func (unspentOutputs *OutputSet) RemoveByTxid(txid string, index uint32) {
	if unspentOutputs.byOutpoint != nil {
		txidBytes, err := hex.DecodeString(txid)
		if err != nil {
			return
		}

		var hash [32]byte
		copy(hash[:], txidBytes)
		common.ReverseBytesInPlace(hash[:])

		outpoint := &tx.PrevOut{Hash: hash, Index: index}
		unspentOutputs.RemoveByOutpoint(outpoint)
	}
}

// UpdateFromBlock mutates the OutputSet by parsing the inputs and outputs of transactions in
// the given block. If unspent outputs in the OutputSet are spent by transactions in the block,
// those outputs are removed from the OutputSet. If new outputs are created which lock funds using
// any of the provided scriptPubKeys, those outputs are added to the set. Transactions are processed
// sequentially, so it is possible for outputs to be created and spent in the same block, which
// returns the OutputSet to its prior state.
//
// This function is very computationally expensive. More
// scriptPubKeys will result in more computational overhead.
func (unspentOutputs *OutputSet) UpdateFromBlock(block *blocks.Block, scriptPubKeys [][]byte) error {
	for _, txn := range block.Transactions {
		for _, vin := range txn.Inputs {
			unspentOutputs.RemoveByOutpoint(vin.PrevOut)
		}

		for outputIndex, vout := range txn.Outputs {
			for _, scriptPubKey := range scriptPubKeys {
				if bytes.Equal(vout.Script, scriptPubKey) {
					txHash, err := txn.Hash(false)
					if err != nil {
						return err
					}

					unspentOutputs.AddOutput(&Output{
						Outpoint: &tx.PrevOut{Hash: txHash, Index: uint32(outputIndex)},
						TxOut: &tx.Output{
							Value:  vout.Value,
							Script: scriptPubKey,
						},
					})
				}
			}
		}
	}

	return nil
}
