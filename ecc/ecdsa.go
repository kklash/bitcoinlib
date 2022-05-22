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
func SignECDSA(privateKey, hash []byte) (r, s *big.Int) {
	d := new(big.Int).SetBytes(privateKey)
	k := Q.Nonce(d, hash, sha256.New)
	z := Q.Bits2int(hash)

	r = new(big.Int)
	s = new(big.Int)
	ekliptic.SignECDSA(d, k, z, r, s)
	return
}

// VerifyECDSA calculates if the given signature (r, s) is a valid ECDSA signature on hash from
// the given public key. Note that non-canonical ECDSA signatures (where s > N/2) are acceptable.
func VerifyECDSA(hash []byte, r, s, pubX, pubY *big.Int) bool {
	// TODO accept public key in bytes form
	z := Q.Bits2int(hash)
	return ekliptic.VerifyECDSA(z, r, s, pubX, pubY)
}
