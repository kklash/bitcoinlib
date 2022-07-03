package ecc

import (
	"math/big"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/common"
	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/ekliptic"
)

var (
	bip340ChallengeHasher = bhash.NewTaggedHasher("BIP0340/challenge")
	bip340AuxHasher       = bhash.NewTaggedHasher("BIP0340/aux")
	bip340NonceHasher     = bhash.NewTaggedHasher("BIP0340/nonce")
)

// SignSchnorr signs a given 32-byte messageHash with the given private key,
// using auxRand as the seed to derive a nonce value.
func SignSchnorr(privateKey, messageHash, auxRand []byte) []byte {
	if len(messageHash) != 32 {
		panic("unexpected message hash length for schnorr signature")
	} else if len(privateKey) != 32 {
		panic("unexpected private key length for schnorr signature")
	} else if len(auxRand) != 32 {
		panic("unexpected aux rand length for schnorr signature")
	}

	d := new(big.Int).SetBytes(privateKey)
	if !IsValidCurveScalar(d) {
		panic("private key is not in range [1, N)")
	}

	pubX := new(big.Int)
	pubY := new(big.Int)
	ekliptic.MultiplyBasePoint(d, pubX, pubY)

	if !isEven(pubY) {
		d.Sub(ekliptic.Secp256k1_CurveOrder, d)
	}

	pubBytes := pubX.FillBytes(make([]byte, 32))

	t := common.XorBytes(
		d.FillBytes(make([]byte, 32)),
		bip340AuxHasher(auxRand),
	)

	rnd := bip340NonceHasher(t, pubBytes, messageHash)

	k := new(big.Int).SetBytes(rnd)
	k.Mod(k, ekliptic.Secp256k1_CurveOrder)
	if equal(k, zero) {
		panic("schnorr signature produced unexpected k of zero")
	}

	rX := new(big.Int)
	rY := new(big.Int)
	ekliptic.MultiplyBasePoint(k, rX, rY)

	if !isEven(rY) {
		k.Sub(ekliptic.Secp256k1_CurveOrder, k)
	}

	rBytes := rX.FillBytes(make([]byte, 32))

	e := new(big.Int).SetBytes(
		bip340ChallengeHasher(
			rBytes,
			pubBytes,
			messageHash,
		),
	)
	e.Mod(e, ekliptic.Secp256k1_CurveOrder)

	s := k.Add(k, e.Mul(e, d))
	s.Mod(s, ekliptic.Secp256k1_CurveOrder)
	k, e = nil, nil

	sig := append(rBytes, s.FillBytes(make([]byte, 32))...)
	return sig
}

// VerifySchnorr returns true if the given signature was made by the owner of the given public key
// on the given message hash.
func VerifySchnorr(pubBytes, messageHash, sig []byte) bool {
	if len(messageHash) != 32 {
		panic("unexpected message hash length for schnorr verification")
	} else if len(sig) != 64 {
		panic("unexpected signature length for schnorr verification")
	}

	if len(pubBytes) != constants.PublicKeySchnorrLength {
		return false
	}

	pubX, pubY, err := DeserializePoint(pubBytes)
	if err != nil {
		return false
	}

	rBytes := sig[:32]
	r := new(big.Int).SetBytes(rBytes)
	if r.Cmp(ekliptic.Secp256k1_P) >= 0 {
		return false
	}

	s := new(big.Int).SetBytes(sig[32:])
	if s.Cmp(ekliptic.Secp256k1_CurveOrder) >= 0 {
		return false
	}

	e := new(big.Int).SetBytes(
		bip340ChallengeHasher(
			r.FillBytes(make([]byte, 32)),
			pubBytes,
			messageHash,
		),
	)
	e.Mod(e, ekliptic.Secp256k1_CurveOrder)

	sgx := new(big.Int)
	sgy := new(big.Int)
	ekliptic.MultiplyBasePoint(s, sgx, sgy)

	epx := new(big.Int)
	epy := new(big.Int)
	ekliptic.MultiplyAffine(pubX, pubY, e, epx, epy, nil)
	ekliptic.Negate(epy)

	Rx := new(big.Int)
	Ry := new(big.Int)
	ekliptic.AddAffine(sgx, sgy, epx, epy, Rx, Ry)

	return !equal(Rx, zero) && !equal(Ry, zero) && isEven(Ry) && equal(Rx, r)
}
