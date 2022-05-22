// Package satutil provides utilities for converting numbers to and from Satoshis and Bitcoins.
package satutil

import (
	"math"
	"math/big"

	"github.com/kklash/bitcoinlib/constants"
)

// SatsToBitcoins converts the given number of Satoshis into Bitcoins.
func SatsToBitcoins(sats uint64) float64 {
	satsBig := new(big.Float).SetUint64(sats)
	satsPerBtc := new(big.Float).SetUint64(constants.SatoshisPerBitcoin)
	btc, _ := satsBig.Quo(satsBig, satsPerBtc).Float64()
	return btc
}

// BitcoinsToSats converts the given number of Bitcoins into Satoshis, rounding
// any significant figures beyond the 8th decimal place.
func BitcoinsToSats(btc float64) uint64 {
	return uint64(math.Round(btc * constants.SatoshisPerBitcoin))
}

// RoundBitcoins rounds the given Bitcoins value to 8 decimal places.
func RoundBitcoins(btc float64) float64 {
	return SatsToBitcoins(BitcoinsToSats(btc))
}
