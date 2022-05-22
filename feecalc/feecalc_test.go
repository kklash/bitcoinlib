package feecalc

import (
	"testing"

	"github.com/kklash/bitcoinlib/tx"
)

func TestFeeCalcTransactions(t *testing.T) {
	for _, fixture := range TestFixtures {
		txn := txFromHex(fixture.TxHex)
		unspentOutputs := fixture.UnspentOutputs

		getPrevOutValue := func(prevOut *tx.PrevOut) (uint64, error) {
			output := unspentOutputs.GetByOutpoint(prevOut)
			if output == nil {
				return 0, ErrPrevOutNotFound
			}

			value := output.TxOut.Value
			return value, nil
		}

		feeValue, err := TotalFeeValue(txn, getPrevOutValue)
		if err != nil {
			t.Errorf("failed to calculate total fee value: %s", err)
			continue
		}

		if feeValue != fixture.FeeValue {
			t.Errorf("fee value was not calculated correctly\nwanted %d\ngot    %d", fixture.FeeValue, feeValue)
			continue
		}

		expectedSatPerByte := float64(fixture.FeeValue) / float64(fixture.VSize)
		satPerByte, _ := FeePerVByte(txn, getPrevOutValue)

		if satPerByte != expectedSatPerByte {
			t.Errorf("sat/vbyte was note calculated correctly\nwanted %f\ngot    %f", expectedSatPerByte, satPerByte)
			continue
		}
	}
}
