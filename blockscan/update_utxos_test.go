package blockscan

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/kklash/bitcoinlib/common"
	"github.com/kklash/bitcoinlib/tx"
	"github.com/kklash/bitcoinlib/unspent"
)

func mustHex(s string) []byte {
	data, _ := hex.DecodeString(s)
	return data
}

func createPrevOut(txid string, vout uint32) *tx.PrevOut {
	prevOut := tx.PrevOut{Index: vout}
	copy(prevOut.Hash[:], common.ReverseBytes(mustHex(txid)))
	return &prevOut
}

func TestBlockScanner_ScanUtxos(t *testing.T) {
	t.Skip()

	scanner, err := createScanner()
	if err != nil {
		t.Errorf("failed to create scanner: %s", err)
		return
	}

	scriptPubKey := mustHex("76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac") // "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx"
	scriptPubKeys := [][]byte{scriptPubKey}
	var startBlock uint32 = 99997
	var endBlock uint32 = 100001

	utxos := new(unspent.OutputSet)
	err = scanner.UpdateUtxos(context.TODO(), utxos, scriptPubKeys, startBlock, endBlock, 6, nil)
	if err != nil {
		t.Errorf("Failed to scan blocks for utxos: %s", err)
		return
	}

	if utxos.Size() != 2 {
		t.Errorf("expected to get output set of size 2; got size %d", utxos.Size())
		return
	}

	expectedPrevOuts := []*tx.PrevOut{
		createPrevOut("fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4", 1),
		createPrevOut("fbde5d03b027d2b9ba4cf5d4fecab9a99864df2637b25ea4cbcb1796ff6550ca", 0),
	}

	for _, prevOut := range expectedPrevOuts {
		output := utxos.GetByOutpoint(prevOut)
		if output == nil {
			t.Errorf("expected to find txid %x", common.ReverseBytes(prevOut.Hash[:]))
			continue
		}

		if output.TxOut.Value != 4444000000 {
			t.Errorf("expected tx value to be 4444000000; got %d", output.TxOut.Value)
		}
	}
}
