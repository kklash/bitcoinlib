package ecc

import (
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"

	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/ekliptic"
)

// ErrPointNotOnCurve is returned upon deserializing an invalid point
// (one which does not satisfy the secp256k1 curve equation).
var ErrPointNotOnCurve = errors.New("failed to deserialize point not on secp256k1 curve")

// DeserializePoint decodes the given serialized curve point, which should either be
// length 65 (uncompressed), 33 (compressed), or 32 (BIP-340 schnorr).
// Returns ErrPointNotOnCurve if the resulting point is not on the secp256k1 curve.
func DeserializePoint(serialized []byte) (x, y *big.Int, err error) {
	switch len(serialized) {
	case constants.PublicKeyUncompressedLength:
		if serialized[0] != constants.PublicKeyUncompressedPrefix {
			return nil, nil, fmt.Errorf("unexpected point prefix byte 0x%x", serialized[0])
		}
		x = new(big.Int).SetBytes(serialized[1:33])
		y = new(big.Int).SetBytes(serialized[33:])

		evenY, oddY := ekliptic.Weierstrass(x)
		if evenY == nil || oddY == nil || !(equal(y, evenY) || equal(y, oddY)) {
			return nil, nil, ErrPointNotOnCurve
		}

	case constants.PublicKeyCompressedLength:
		x = new(big.Int).SetBytes(serialized[1:])
		evenY, oddY := ekliptic.Weierstrass(x)
		if evenY == nil || oddY == nil {
			return nil, nil, ErrPointNotOnCurve
		}

		switch serialized[0] {
		case constants.PublicKeyCompressedEvenByte:
			y = evenY
		case constants.PublicKeyCompressedOddByte:
			y = oddY
		default:
			return nil, nil, fmt.Errorf("unexpected point prefix byte 0x%x", serialized[0])
		}

	case constants.PublicKeySchnorrLength:
		x = new(big.Int).SetBytes(serialized)
		y, _ = ekliptic.Weierstrass(x)
		if y == nil {
			return nil, nil, ErrPointNotOnCurve
		}

	default:
		err = fmt.Errorf("attempted to deserialize unexpected byte slice of length %d as curve point", len(serialized))
	}

	return
}

// SerializePointUncompressed serializes the given curve point in 65-byte uncompressed format.
func SerializePointUncompressed(x, y *big.Int) []byte {
	return elliptic.Marshal(Curve, x, y)
}

// SerializePointCompressed serializes the given curve point in 33-byte compressed format.
func SerializePointCompressed(x, y *big.Int) []byte {
	return elliptic.MarshalCompressed(Curve, x, y)
}

// SerializePoint serializes the given curve point in compressed format if compressed is true,
// otherwise it returns the uncompressed serialization.
func SerializePoint(x, y *big.Int, compressed bool) []byte {
	if compressed {
		return SerializePointCompressed(x, y)
	}
	return SerializePointUncompressed(x, y)
}
