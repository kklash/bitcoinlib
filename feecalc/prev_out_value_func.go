package feecalc

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/kklash/bitcoinlib/common"
	"github.com/kklash/bitcoinlib/tx"
)

// ErrPrevOutNotFound should be returned (or wrapped and returned) by a PrevOutValueFunc
// when the requested PrevOut cannot be found.
var ErrPrevOutNotFound = errors.New("requested prev out could not be found")

// A PrevOutValueFunc should return the value of a given previous output in satoshis.
// Should return ErrPrevOutNotFound if the given PrevOut could not be found.
type PrevOutValueFunc func(*tx.PrevOut) (uint64, error)

// NewNaivePrevOutValueFunc returns a naive non-caching implementation of a PrevOutValueFunc
// using transaction data fetched by the given getTxHex function. Each call to the
// PrevOutValueFunc fetches the transaction involved in the requested previous output.
// The transaction is then parsed and the output value returned.
func NewNaivePrevOutValueFunc(getTxHex func(txid string) (string, error)) PrevOutValueFunc {
	return func(prevOut *tx.PrevOut) (uint64, error) {
		txid := hex.EncodeToString(common.ReverseBytes(prevOut.Hash[:]))

		txHex, err := getTxHex(txid)
		if err != nil {
			return 0, err
		} else if txHex == "" {
			return 0, fmt.Errorf("%w: txid not found - %s", ErrPrevOutNotFound, txid)
		}

		txBytes, err := hex.DecodeString(txHex)
		if err != nil {
			return 0, err
		}

		txn, err := tx.FromBytes(txBytes)
		if err != nil {
			return 0, err
		}

		if int(prevOut.Index) >= len(txn.Outputs) {
			return 0, fmt.Errorf("%w: bad index - %s:%d", ErrPrevOutNotFound, txid, prevOut.Index)
		}

		value := txn.Outputs[prevOut.Index].Value
		return value, nil
	}
}
