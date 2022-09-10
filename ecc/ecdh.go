package ecc

import (
	"math/big"

	"github.com/kklash/ekliptic"
)

// SharedSecret generates a shared secret based on a private key and a
// public key using Diffie-Hellman key exchange (ECDH) (RFC 4753).
// RFC5903 Section 9 states we should only return x.
func SharedSecret(priv, pubX, pubY *big.Int) []byte {
	sharedKey, _ := ekliptic.MultiplyAffine(pubX, pubY, priv, nil)
	return sharedKey.FillBytes(make([]byte, 32))
}
