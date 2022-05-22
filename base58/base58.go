// Package base58 provides Base 58 byte encoding using the Bitcoin standard base 58 alphabet.
package base58

import (
	"errors"
	"math/big"

	"github.com/kklash/bitcoinlib/constants"
)

// Alphabet is the Bitcoin Base 58 alphabet.
const Alphabet = constants.Base58Alphabet

var (
	radix = big.NewInt(int64(len(Alphabet)))
	zero  = big.NewInt(0)

	// AlphabetIndeces is a mapping of every char in Alphabet to their index in Alphabet.
	AlphabetIndeces = make(map[byte]int64)

	// ErrInvalidBase58String is returned by Decode when it is passed an improperly encoded base 58 string.
	ErrInvalidBase58String = errors.New("Unable to decode invalid base 58 string")
)

func init() {
	for i := 0; i < len(Alphabet); i++ {
		r := Alphabet[i]
		AlphabetIndeces[r] = int64(i)
	}
}

// Encode encodes data into a base 58 string and returns the result.
func Encode(data []byte) (bs58 string) {
	if data == nil {
		return
	}

	bi := new(big.Int).SetBytes(data)

	for bi.Cmp(zero) != 0 {
		_, rem := bi.QuoRem(bi, radix, new(big.Int))
		bs58 = string(Alphabet[rem.Int64()]) + bs58
	}

	// For number of leading 0's in bytes, prepend 1
	for _, b := range data {
		if b == 0 {
			bs58 = string(Alphabet[0]) + bs58
		} else {
			break
		}
	}

	return
}

// Decode decodes a base 58 string and returns the resulting bytes. Returns ErrInvalidBase58String
// If the string bs58 is not properly encoded.
func Decode(bs58 string) ([]byte, error) {
	nZeros := 0
	for nZeros = 0; nZeros < len(bs58) && bs58[nZeros] == Alphabet[0]; nZeros++ {
	}

	bs58 = bs58[nZeros:]

	ret := big.NewInt(0)

	for i := 0; i < len(bs58); i++ {
		c := bs58[len(bs58)-i-1]
		m, ok := AlphabetIndeces[c]
		if !ok {
			return nil, ErrInvalidBase58String
		}

		// base58 equiv of ret += m << i
		expVal := new(big.Int).Exp(radix, big.NewInt(int64(i)), nil)
		expVal.Mul(big.NewInt(m), expVal)

		ret.Add(ret, expVal)
	}

	decoded := append(make([]byte, nZeros), ret.Bytes()...)

	return decoded, nil
}
