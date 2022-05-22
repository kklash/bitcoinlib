// Package ecc provides an implementation of elliptic curves of short Weierstass form
// over finite fields, used for cryptography.
package ecc

import (
	"math/big"
)

var (
	zero  = big.NewInt(0)
	seven = big.NewInt(7)
)

func equal(v1, v2 *big.Int) bool {
	if v1 == nil || v2 == nil {
		return v1 == v2
	}

	return v1.Cmp(v2) == 0
}

func bigIntFromHex(s string) *big.Int {
	i, _ := new(big.Int).SetString(s, 16)
	return i
}

func pointAt(x, y int64) (*big.Int, *big.Int) {
	return big.NewInt(x), big.NewInt(y)
}

func pointAtHex(xHex, yHex string) (*big.Int, *big.Int) {
	x := new(big.Int)
	x.SetString(xHex, 16)
	y := new(big.Int)
	y.SetString(yHex, 16)
	return x, y
}

func isEven(y *big.Int) bool {
	return y.Bit(0) == 0
}
