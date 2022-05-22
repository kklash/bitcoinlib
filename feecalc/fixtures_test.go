package feecalc

import (
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/kklash/bitcoinlib/common"
	"github.com/kklash/bitcoinlib/tx"
	"github.com/kklash/bitcoinlib/unspent"
)

type Fixture struct {
	// Directly from JSON
	TxHex       string `json:"txHex"`
	FeeValue    uint64 `json:"feeValue"`
	VSize       int    `json:"vsize"`
	PrevOutsRaw []*struct {
		Txid  string `json:"txid"`
		Vout  uint32 `json:"vout"`
		Value uint64 `json:"value"`
	} `json:"prevOuts"`

	// Reformatted from json
	Tx             *tx.Tx
	UnspentOutputs *unspent.OutputSet
}

var TestFixtures []*Fixture

func createUnspentOutput(txid string, vout uint32, value uint64) *unspent.Output {
	prevOut := tx.PrevOut{Index: vout}
	txidBytes, _ := hex.DecodeString(txid)
	copy(prevOut.Hash[:], common.ReverseBytes(txidBytes))
	return &unspent.Output{
		Outpoint: &prevOut,
		TxOut:    &tx.Output{Value: value},
	}
}

func txFromHex(txHex string) *tx.Tx {
	txBytes, _ := hex.DecodeString(txHex)
	txn, _ := tx.FromBytes(txBytes)
	return txn
}

func init() {
	fh, err := os.Open("fixtures.json")
	if err != nil {
		panic(err)
	}

	if err := json.NewDecoder(fh).Decode(&TestFixtures); err != nil {
		panic(err)
	}

	for _, fixture := range TestFixtures {
		outputs := make([]*unspent.Output, len(fixture.PrevOutsRaw))
		for i, prevOutRaw := range fixture.PrevOutsRaw {
			outputs[i] = createUnspentOutput(prevOutRaw.Txid, prevOutRaw.Vout, prevOutRaw.Value)
		}

		fixture.UnspentOutputs = unspent.NewOutputSet(outputs)
		fixture.Tx = txFromHex(fixture.TxHex)
	}
}
