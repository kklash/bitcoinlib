package unspent

import (
	"testing"

	"github.com/kklash/bitcoinlib/tx"
)

func TestOutputSet(t *testing.T) {
	unspentOutputs := new(OutputSet)

	if unspentOutputs.Size() != 0 {
		t.Errorf("expected empty OutputSet size to be zero, got %d", unspentOutputs.Size())
		return
	}

	utxo1 := &Output{
		Outpoint: &tx.PrevOut{
			Hash:  mustHash("0000000000000000000000000000000000000000000000000000000000000001"),
			Index: 1,
		},
		TxOut: &tx.Output{
			Value:  20410,
			Script: mustHex("001481916950b977370407e58cb3970ce1292093e6e6"),
		},
	}

	unspentOutputs.AddOutput(utxo1)

	if unspentOutputs.Size() != 1 {
		t.Errorf("expected OutputSet size to be one, got %d", unspentOutputs.Size())
		return
	}

	gotUtxo := unspentOutputs.GetByHash(utxo1.Outpoint.Hash, utxo1.Outpoint.Index)
	if gotUtxo != utxo1 {
		t.Errorf("expected to retrieve utxo by hash and output index")
		return
	}

	gotUtxo = unspentOutputs.GetByTxid("0100000000000000000000000000000000000000000000000000000000000000", utxo1.Outpoint.Index)
	if gotUtxo != utxo1 {
		t.Errorf("expected to retrieve utxo by hash and output index")
		return
	}

	unspentOutputs.AddOutput(utxo1)
	if unspentOutputs.Size() != 1 {
		t.Errorf("expected OutputSet size to be one, got %d", unspentOutputs.Size())
		return
	}

	utxo2 := &Output{
		Outpoint: &tx.PrevOut{
			Hash:  mustHash("0000000000000000000000000000000000000000000000000000000000000002"),
			Index: 2,
		},
		TxOut: &tx.Output{
			Value:  30000,
			Script: mustHex("76a914d32a2b0ff8a8eb2471157217d25a11c05b072cfb88ac"),
		},
	}
	unspentOutputs.AddOutput(utxo2)

	slice := unspentOutputs.Slice()
	if len(slice) != 2 ||
		(slice[0] != utxo1 && slice[0] != utxo2) ||
		(slice[1] != utxo1 && slice[1] != utxo2) {
		t.Errorf("failed to derive slice from OutputSet")
		return
	}

	if unspentOutputs.Size() != 2 {
		t.Errorf("expected OutputSet size to be two, got %d", unspentOutputs.Size())
		return
	}

	unspentOutputs.RemoveByHash(utxo2.Outpoint.Hash, utxo2.Outpoint.Index)
	if unspentOutputs.Size() != 1 {
		t.Errorf("expected OutputSet size to be one, got %d", unspentOutputs.Size())
		return
	}
}
