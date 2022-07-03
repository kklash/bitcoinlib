package ecc

import (
	"crypto/sha256"
	"math/big"

	"github.com/kklash/ekliptic"
	"github.com/kklash/rfc6979"
)

var Q = rfc6979.NewQ(ekliptic.Secp256k1_CurveOrder)

// SignECDSA signs the given hash with the given private key d and
// returns the two components of the signature: r and s.
//
// SignECDSA calculates the secret signature nonce value k deterministically using RFC6979.
func SignECDSA(privateKey, messageHash []byte) (r, s *big.Int) {
	if len(messageHash) != 32 {
		panic("unexpected message hash length for ECDSA signature")
	} else if len(privateKey) != 32 {
		panic("unexpected private key length for ECDSA signature")
	}

	d := new(big.Int).SetBytes(privateKey)
	k := Q.Nonce(d, messageHash, sha256.New)
	z := Q.Bits2int(messageHash)

	r = new(big.Int)
	s = new(big.Int)
	ekliptic.SignECDSA(d, k, z, r, s)
	return
}

// VerifyECDSA calculates if the given signature (r, s) is a valid ECDSA signature on messageHash from
// the given public key. Note that non-canonical ECDSA signatures (where s > N/2) are acceptable.
func VerifyECDSA(pubBytes, messageHash []byte, r, s *big.Int) bool {
	if len(messageHash) != 32 {
		panic("unexpected message hash length for ECDSA verification")
	}

	if !IsValidCurveScalar(r) || !IsValidCurveScalar(s) {
		return false
	}

	pubX, pubY, err := DeserializePoint(pubBytes)
	if err != nil {
		return false
	}

	z := Q.Bits2int(messageHash)
	return ekliptic.VerifyECDSA(z, r, s, pubX, pubY)
}
