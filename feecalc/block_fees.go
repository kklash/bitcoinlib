package feecalc

import (
	"github.com/kklash/bitcoinlib/blocks"
)

// FeeRangeForBlock returns the min and max fee rates for transactions found in the given block's transactions.
func FeeRangeForBlock(block *blocks.Block, getPrevOutValue PrevOutValueFunc) (min, max float64, err error) {
	min = -1.0
	max = 0.0

	for _, txn := range block.Transactions[1:] {
		satoshisPerVByte, err := FeePerVByte(txn, getPrevOutValue)
		if err != nil {
			return 0, 0, err
		}

		if satoshisPerVByte > max {
			max = satoshisPerVByte
		}
		if satoshisPerVByte < min || min == -1.0 {
			min = satoshisPerVByte
		}
	}

	// In case block is only a coinbase transaction
	if min == -1.0 {
		min = 0.0
	}

	return
}

// TotalFeesForBlock returns the total amount of satoshis paid as miner
// fees for all transactions in the given block.
func TotalFeesForBlock(block *blocks.Block, getPrevOutValue PrevOutValueFunc) (uint64, error) {
	var total uint64 = 0

	for _, txn := range block.Transactions[1:] {
		feeForTxn, err := TotalFeeValue(txn, getPrevOutValue)
		if err != nil {
			return 0, err
		}

		total += feeForTxn
	}

	return total, nil
}

// AverageFeeForBlockPerVByte returns the average fee paid per virtual byte of block space.
// This is the total fees paid in the block divided by the block weight, divided again by four.
func AverageFeeForBlockPerVByte(block *blocks.Block, getPrevOutValue PrevOutValueFunc) (float64, error) {
	totalFeesForBlock, err := TotalFeesForBlock(block, getPrevOutValue)
	if err != nil {
		return 0, err
	}

	blockWeight := block.WeightUnits()

	averageFeePerWeightUnit := float64(totalFeesForBlock) / float64(blockWeight)
	averageFeePerVByte := averageFeePerWeightUnit / 4

	return averageFeePerVByte, nil
}
