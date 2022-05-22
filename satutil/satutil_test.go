package satutil

import "testing"

func TestSatoshis(t *testing.T) {
	type Fixture struct {
		btc        float64
		roundedBtc float64
		sats       uint64
	}

	fixtures := []*Fixture{
		{0.00000001, 0.00000001, 1},
		{1.0, 1.0, 100_000_000},
		{1.43128382, 1.43128382, 143_128_382},
		{173.1, 173.1, 17_310_000_000},
		{1.000000009, 1.00000001, 100_000_001},
		{12_000_000.173812819, 12_000_000.17381282, 1_200_000_017_381_282},
	}

	for _, fixture := range fixtures {
		sats := BitcoinsToSats(fixture.btc)
		if sats != fixture.sats {
			t.Errorf("failed to convert BTC to Satoshis\nWanted %d\nGot    %d", fixture.sats, sats)
			continue
		}

		roundedBtc := RoundBitcoins(fixture.btc)
		if roundedBtc != fixture.roundedBtc {
			t.Errorf("failed to round BTC to 8 decimals\nWanted %v\nGot    %v", fixture.roundedBtc, roundedBtc)
			continue
		}

		btc := SatsToBitcoins(sats)
		if btc != roundedBtc {
			t.Errorf("failed to convert Satoshis to BTC\nWanted %v\nGot    %v", roundedBtc, btc)
			continue
		}
	}
}
