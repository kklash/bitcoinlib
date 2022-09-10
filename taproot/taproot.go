package taproot

import (
	"fmt"
	"math/big"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/ecc"
	"github.com/kklash/ekliptic"
)

var tapTweakHasher = bhash.NewTaggedHasher("TapTweak")

// TweakPublicKey derives a public key that commits to the given public key
// and the taproot commitment h.
func TweakPublicKey(publicKey, h []byte) (qPub []byte, hasOddY bool, err error) {
	if len(publicKey) != 32 {
		err = fmt.Errorf("expected 32-byte schnorr public key; got %d bytes", len(publicKey))
		return
	}
	pubX, pubY, err := ecc.DeserializePoint(publicKey)
	if err != nil {
		return
	}

	t := new(big.Int).SetBytes(tapTweakHasher(append(publicKey, h...)))
	if !ekliptic.IsValidScalar(t) {
		err = fmt.Errorf("invalid tweaked public key; t exceeds curve order")
		return
	}

	tGx, tGy := ekliptic.MultiplyBasePoint(t)

	// Q = P + hash(P, h)G
	qx, qy := ekliptic.AddAffine(pubX, pubY, tGx, tGy)

	qPub = qx.FillBytes(make([]byte, 32))
	hasOddY = qy.Bit(0) == 1
	return
}

// TweakPrivateKey derives a private key that commits to the given private key and
// the taproot commitment h.
//
// For a given private/public key pair, and any commitment value h, it holds that
// the private key tweaked with h controls the public key tweaked with h.
func TweakPrivateKey(privateKey, h []byte) ([]byte, error) {
	seckey := new(big.Int).SetBytes(privateKey)
	if !ekliptic.IsValidScalar(seckey) {
		return nil, fmt.Errorf("cannot tweak invalid private key")
	}

	pubX, pubY := ekliptic.MultiplyBasePoint(seckey)

	if pubY.Bit(0) == 1 {
		seckey = new(big.Int).Sub(ekliptic.Secp256k1_CurveOrder, seckey)
	}

	t := new(big.Int).SetBytes(tapTweakHasher(append(pubX.FillBytes(make([]byte, 32)), h...)))
	if !ekliptic.IsValidScalar(t) {
		return nil, fmt.Errorf("invalid tweaked private key; t exceeds curve order")
	}
	seckey.Add(seckey, t)
	seckey.Mod(seckey, ekliptic.Secp256k1_CurveOrder)
	tweakedPriv := seckey.FillBytes(make([]byte, 32))
	return tweakedPriv, nil
}
