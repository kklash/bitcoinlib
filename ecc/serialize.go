package ecc

import (
	"fmt"
	"math/big"

	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/ekliptic"
)

func Serialize256(v *big.Int) []byte {
	buf := make([]byte, 32)
	v.FillBytes(buf)
	return buf
}

func SerializePoint(x, y *big.Int) []byte {
	buf := make([]byte, 64)
	x.FillBytes(buf[:32])
	y.FillBytes(buf[32:])
	return buf
}

func SerializePointCompressed(x, y *big.Int) []byte {
	var prefix byte
	if isEven(y) {
		prefix = constants.PublicKeyCompressedEvenByte // Even number
	} else {
		prefix = constants.PublicKeyCompressedOddByte // Odd number
	}

	buf := make([]byte, 33)
	buf[0] = prefix
	x.FillBytes(buf[1:])
	return buf
}

// TODO return an error instead of panic
func DeserializePoint(serialized []byte) (x, y *big.Int) {
	switch len(serialized) {
	case constants.PublicKeyUncompressedLength:
		x = new(big.Int).SetBytes(serialized[:32])
		y = new(big.Int).SetBytes(serialized[32:])

	case constants.PublicKeyCompressedLength:
		x = new(big.Int).SetBytes(serialized[1:])
		evenY, oddY := ekliptic.Weierstrass(x)
		if serialized[0]%2 == 0 {
			y = evenY
		} else {
			y = oddY
		}

	// TODO 32-byte x-only keys from bip340

	default:
		panic(fmt.Sprintf("attempted to deserialize unexpected byte slice of length %d", len(serialized)))
	}

	if !ekliptic.IsOnCurveAffine(x, y) {
		panic("deserialized public key not on the secp256k1 curve")
	}

	return
}
