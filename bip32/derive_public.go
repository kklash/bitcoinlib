package bip32

import (
	"fmt"

	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/bitcoinlib/ecc"
)

// ErrInvalidDerivationIndex is returned by DerivePublicChild
// if the given index is above the hardening threshold.
var ErrInvalidDerivationIndex = fmt.Errorf("Cannot derive child of public key with index higher than 0x%x", constants.Bip32Hardened)

func derivePublicChild(parentPublicKey, chainCode []byte, childIndex uint32) (childPublicKey, childChainCode []byte, err error) {
	if childIndex >= constants.Bip32Hardened {
		err = ErrInvalidDerivationIndex
		return
	}

	pubX, pubY, err := ecc.DeserializePoint(parentPublicKey)
	if err != nil {
		return
	}

	data := append(parentPublicKey, serialize32(childIndex)...)
	l := hmacSha512(chainCode, data)
	lLeft, lRight := l[:32], l[32:]

	lLeftGx, lLeftGy := curve.ScalarBaseMult(lLeft)

	childKeyX, childKeyY := curve.Add(lLeftGx, lLeftGy, pubX, pubY)

	childPublicKey = ecc.SerializePointCompressed(childKeyX, childKeyY)
	childChainCode = lRight
	return
}

// DerivePublicChild derives a child public key and chain code of a given
// bip32 path from a parent public key and chain code. Returns ErrInvalidDerivationIndex
// if any child index is greater than or equal to constants.Bip32Hardened.
//
// The parentPublicKey parameter can be compressed or uncompressed, but this function
// always returns compressed public keys.
func DerivePublicChild(parentPublicKey, chainCode []byte, childIndices ...uint32) ([]byte, []byte, error) {
	if len(childIndices) == 0 {
		return parentPublicKey, chainCode, nil
	}

	newKey, newChainCode, err := derivePublicChild(parentPublicKey, chainCode, childIndices[0])
	if err != nil {
		return nil, nil, err
	}

	return DerivePublicChild(newKey, newChainCode, childIndices[1:]...)
}
