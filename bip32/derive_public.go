package bip32

import (
	"fmt"
	"math/big"

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
	} else if err = ValidatePublicKeyBytes(parentPublicKey); err != nil {
		return
	}

	data := append(parentPublicKey, serialize32(childIndex)...)
	l := hmacSha512(chainCode, data)
	lLeft, lRight := l[:32], l[32:]

	lLeftGx, lLeftGy := curve.ScalarBaseMult(lLeft)
	var pubX, pubY *big.Int

	if IsCompressedPublicKey(parentPublicKey) {
		if pubX, pubY, err = UncompressPublicKey(parentPublicKey); err != nil {
			return
		}
	} else {
		xBytes, yBytes := parentPublicKey[:32], parentPublicKey[32:]
		pubX, pubY = parse256(xBytes), parse256(yBytes)
	}

	childKeyX, childKeyY := curve.Add(lLeftGx, lLeftGy, pubX, pubY)

	childPublicKey = ecc.SerializePointCompressed(childKeyX, childKeyY)
	childChainCode = lRight
	return
}

// DerivePublicChild derives a child public key and chain code of a given
// bip32 path from a parent public key and chain code. Returns ErrInvalidDerivationIndex
// if any child index is greater than or equal to constants.Bip32Hardened. Returns ErrInvalidPublicKey
// if parentPublicKey is of an invalid format, or not on the secp256k1 curve. The
// parentPublicKey parameter can be compressed or uncompressed, but this function
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
