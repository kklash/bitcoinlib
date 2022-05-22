package unspent

import (
	"github.com/kklash/bitcoinlib/tx"
)

// Output represents an unspent transaction output, including a
// reference to the transaction outpoint which created it.
type Output struct {
	Outpoint *tx.PrevOut
	TxOut    *tx.Output
}

// Clone returns a pointer to a duplicate of the unspent Output.
func (o *Output) Clone() *Output {
	clone := &Output{
		Outpoint: o.Outpoint.Clone(),
		TxOut:    o.TxOut.Clone(),
	}

	return clone
}
