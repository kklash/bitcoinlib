package blockheader

import (
	"math/big"
)

// https://developer.bitcoin.org/reference/block_chain.html#target-nbits
func calculateTargetNBits(nBits uint32) *big.Int {
	significand := nBits & 0x00ffffff

	// Handle negative bit set in NBits significand.
	// When parsing nBits, Bitcoin Core converts a negative target threshold into a
	// target of zero, which the header hash can equal (in theory, at least).
	if significand&0x800000 != 0 {
		return big.NewInt(0)
	}

	exponent := nBits >> (8 * 3)

	// for small exponents values, multiplying a bignum by a negative exponent will give you 1 instead of zero.
	// so instead we right shift the significand to discard the unused bytes.
	if exponent < 3 {
		significand >>= (8 * (3 - exponent))
	} else if exponent > 32 {
		panic("cannot calculate target nbits; exponent overflows uint256")
	}

	bigSignificand := big.NewInt(int64(significand))
	bigExponent := big.NewInt(int64(exponent) - 3)
	base := big.NewInt(256)
	basePowExponent := new(big.Int).Exp(base, bigExponent, nil)
	target := new(big.Int).Mul(bigSignificand, basePowExponent)
	return target
}
