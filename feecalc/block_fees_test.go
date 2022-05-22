package feecalc

// TODO redo this testing completely

// func TestFeeRangeForBlock(t *testing.T) {
// 	type BlockFixture struct {
// 		block      *blocks.Block
// 		totalFees  uint64
// 		minFeeRate float64
// 		maxFeeRate float64
// 	}

// 	createBlockFixture := func(txFixtures ...*Fixture) *BlockFixture {
// 		blockFixture := &BlockFixture{
// 			block: &blocks.Block{
// 				Transactions: make([]*tx.Tx, len(txFixtures)+1),
// 			},
// 			totalFees:  0,
// 			minFeeRate: float64(txFixtures[0].FeeValue) / float64(txFixtures[0].VSize),
// 			maxFeeRate: float64(txFixtures[0].FeeValue) / float64(txFixtures[0].VSize),
// 		}

// 		// coinbase
// 		blockFixture.block.Transactions[0] = TestFixtures[0].Tx

// 		for i, fixture := range txFixtures {
// 			feeRate := float64(fixture.FeeValue) / float64(fixture.VSize)
// 			if feeRate < blockFixture.minFeeRate {
// 				blockFixture.minFeeRate = feeRate
// 			}
// 			if feeRate > blockFixture.maxFeeRate {
// 				blockFixture.maxFeeRate = feeRate
// 			}

// 			blockFixture.block.Transactions[i+1] = fixture.Tx
// 			blockFixture.totalFees += fixture.FeeValue
// 		}

// 		return blockFixture
// 	}

// 	blockFixtures := []*BlockFixture{
// 		createBlockFixture(),
// 		createBlockFixture(TestFixtures[0]),
// 		createBlockFixture(TestFixtures[0], TestFixtures[1]),
// 		createBlockFixture(TestFixtures[1], TestFixtures[2]),
// 		createBlockFixture(TestFixtures[2], TestFixtures[3], TestFixtures[4]),
// 		createBlockFixture(TestFixtures[2], TestFixtures[5]),
// 	}

// 	getPrevOutValue := func(prevOut *input.PrevOut) (uint64, error) {
// 		for _, fixture := range TestFixtures {
// 			output := fixture.UnspentOutputs.GetByOutpoint(prevOut)
// 			if output != nil {
// 				return output.TxOut.Value, nil
// 			}
// 		}
// 		return 0, ErrPrevOutNotFound
// 	}

// 	for _, fixture := range blockFixtures {
// 		totalFeesForBlock, err := TotalFeesForBlock(fixture.block, getPrevOutValue)
// 		if err != nil {
// 			t.Errorf("failed to get total block fees: %s", err)
// 			continue
// 		}

// 		if totalFeesForBlock != fixture.totalFees {
// 			t.Errorf("totalFees does not match expected for block\nWanted %d\nGot    %d", fixture.totalFees, totalFeesForBlock)
// 			continue
// 		}

// 		minFeeRate, maxFeeRate, err := FeeRangeForBlock(fixture.block, getPrevOutValue)
// 		if err != nil {
// 			t.Errorf("failed to get fee rate range: %s", err)
// 			continue
// 		}

// 		if minFeeRate != fixture.minFeeRate {
// 			t.Errorf("minFeeRate does not match expected\nwanted %f\ngot    %f", fixture.minFeeRate, minFeeRate)
// 			continue
// 		}
// 		if maxFeeRate != fixture.maxFeeRate {
// 			t.Errorf("maxFeeRate does not match expected\nwanted %f\ngot    %f", fixture.maxFeeRate, maxFeeRate)
// 			continue
// 		}

// 		weight := fixture.block.WeightUnits()
// 		expectedAverageFee := float64(weight) / float64(fixture.totalFees) / 4.0

// 		averageFee, err := AverageFeeForBlockPerVByte(fixture.block, getPrevOutValue)
// 		if err != nil {
// 			t.Errorf("failed to get fee rate average: %s", err)
// 			continue
// 		}

// 		if averageFee != expectedAverageFee {
// 			t.Errorf("average fee does not match expected\nWanted %f\nGot    %f", expectedAverageFee, averageFee)
// 			continue
// 		}
// 	}
// }
