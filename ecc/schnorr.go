package ecc

import (
	"crypto/sha256"
	"math/big"

	"github.com/kklash/ekliptic"
)

func newTaggedHasher(tag string) func(...[]byte) []byte {
	hashedTag := sha256.Sum256([]byte(tag))
	return func(chunks ...[]byte) []byte {
		h := sha256.New()
		h.Write(append(hashedTag[:], hashedTag[:]...))
		for _, chunk := range chunks {
			h.Write(chunk)
		}
		return h.Sum(nil)
	}
}

var (
	bip340ChallengeHasher = newTaggedHasher("BIP0340/challenge")
	bip340AuxHasher       = newTaggedHasher("BIP0340/aux")
	bip340NonceHasher     = newTaggedHasher("BIP0340/nonce")
)

func liftX(x *big.Int) *big.Int {
	y, _ := ekliptic.Weierstrass(x)
	return y
}

// func (curve *Curve) schnorrSign(d *big.Int, hash, auxRand []byte) (*big.Int, *big.Int) {
// 	if len(hash) != 32 {
// 		panic("unexpected hash length for schnorr signature")
// 	} else if len(auxRand) != 32 {
// 		panic("unexpected aux rand data length for schnorr signature")
// 	}

// 	// k := curve.q.Nonce(d, hash, sha256.New)
// 	a := parse256(auxRand)
// 	rX, _ := curve.multiplyScalar(curve.Gx, curve.Gy, a)
// 	pubX, pubY := curve.multiplyScalar(curve.Gx, curve.Gy, d)

// 	if !isEven(pubY) {
// 		d = new(big.Int).Sub(curve.N, d)
// 	}

// 	var hashable [96]byte
// 	copy(hashable[:], Serialize256(rX))
// 	copy(hashable[32:], Serialize256(pubX))
// 	copy(hashable[64:], hash)

// 	hashed := sha256.Sum256(hashable[:])
// 	hashedInt := new(big.Int).SetBytes(hashed[:])
// 	hd := new(big.Int).Mul(hashedInt, d)

// 	// s = k - d * SHA(Rx || pkX || h)
// 	s := new(big.Int).Sub(k, hd)

// 	// TODO additional changes needed, read full spec
// 	// https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki
// 	return rX, s
// }
