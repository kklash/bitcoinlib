// Package feecalc provides fee calculation and estimation utilities.
package feecalc

import (
	"errors"

	"github.com/kklash/bitcoinlib/tx"
)

// ErrInvalidFeeRate is returned when attempting to calculate fee values for
// a transaction whose outputs' value is greater than its inputs' value.
var ErrInvalidFeeRate = errors.New("input value for TX was less than output value")

// TotalOutputValue returns the total amount spent by a transaction's
// outputs in satoshis. Note that this does not include mining fees.
func TotalOutputValue(txn *tx.Tx) uint64 {
	var total uint64 = 0
	for _, vout := range txn.Outputs {
		total += vout.Value
	}

	return total
}

// TotalInputValue returns the total amount spent by a transaction including miner's fees.
// This is the sum of all input values. Must pass a PrevOutValueFunc which provides a source
// of truth for the value of the previous outputs which are being spent by the txn.
func TotalInputValue(txn *tx.Tx, getPrevOutValue PrevOutValueFunc) (uint64, error) {
	var total uint64 = 0
	for _, vin := range txn.Inputs {
		value, err := getPrevOutValue(vin.PrevOut)
		if err != nil {
			return 0, err
		}

		total += value
	}

	return total, nil
}

// TotalFeeValue returns the total fee value of a transaction: the sum of its inputs'
// values minus the sum of its outputs' values. Must pass a PrevOutValueFunc which provides
// a source of truth for the value of the previous outputs which are being spent by the txn.
// Returns ErrInvalidFeeRate if it encounters a transaction whose outputs spend more than its inputs.
func TotalFeeValue(txn *tx.Tx, getPrevOutValue PrevOutValueFunc) (uint64, error) {
	outputValue := TotalOutputValue(txn)
	inputValue, err := TotalInputValue(txn, getPrevOutValue)
	if err != nil {
		return 0, err
	}

	if inputValue < outputValue {
		return 0, ErrInvalidFeeRate
	}

	feeValue := inputValue - outputValue
	return feeValue, nil
}

// FeePerVByte returns the fee rate per byte of the transaction's
// virtual size. A byte of witness data counts for 1/4 of a virtual byte.
// Returns ErrInvalidFeeRate if it encounters a transaction whose outputs spend more than its inputs.
func FeePerVByte(txn *tx.Tx, getPrevOutValue PrevOutValueFunc) (float64, error) {
	feeValue, err := TotalFeeValue(txn, getPrevOutValue)
	if err != nil {
		return 0, err
	}

	satoshisPerVByte := float64(feeValue) / float64(txn.VSize())
	// fmt.Printf("%f / %f = %f\n", float64(feeValue), float64(txn.VSize()), satoshisPerVByte)
	return satoshisPerVByte, nil
}
