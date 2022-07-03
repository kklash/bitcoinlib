package bip32

import (
	"math/big"

	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/bitcoinlib/ecc"
)

func derivePrivateChild(parentPrivateKey, chainCode []byte, childIndex uint32) (childPrivateKey, childChainCode []byte) {
	data := make([]byte, 0)
	if childIndex >= constants.Bip32Hardened {
		data = append(data, 0)
		data = append(data, parentPrivateKey...)
	} else {
		pubX, pubY := curve.ScalarBaseMult(parentPrivateKey)
		data = append(data, ecc.SerializePointCompressed(pubX, pubY)...)
	}
	data = append(data, serialize32(childIndex)...)

	l := hmacSha512(chainCode, data)
	lLeft, lRight := l[:32], l[32:]

	childKeyInt := new(big.Int).SetBytes(lLeft)
	childKeyInt.Add(childKeyInt, new(big.Int).SetBytes(parentPrivateKey))
	childKeyInt.Mod(childKeyInt, curve.Params().N)

	childPrivateKey = childKeyInt.FillBytes(make([]byte, 32))
	childChainCode = lRight
	return
}

// DerivePrivateChild recursively derives a child private key and chain code
// of a given bip32 path from a parent private key and chain code.
func DerivePrivateChild(parentPrivateKey, chainCode []byte, childIndices ...uint32) ([]byte, []byte) {
	if len(childIndices) == 0 {
		return parentPrivateKey, chainCode
	}

	newKey, newChainCode := derivePrivateChild(parentPrivateKey, chainCode, childIndices[0])
	return DerivePrivateChild(newKey, newChainCode, childIndices[1:]...)
}
