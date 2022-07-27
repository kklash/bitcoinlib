package taproot_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/kklash/bitcoinlib/ecc"
	"github.com/kklash/bitcoinlib/taproot"
)

func TestTweak(t *testing.T) {
	privateKey, err := ecc.NewPrivateKey(rand.Reader)
	if err != nil {
		t.Errorf("failed to generate private key: %s", err)
		return
	}
	publicKey := ecc.GetPublicKeySchnorr(privateKey)

	commitment := make([]byte, 32)

	tweakedPriv, err := taproot.TweakPrivateKey(privateKey, commitment)
	if err != nil {
		t.Errorf("failed to tweak private key: %s", err)
		return
	}

	tweakedPub, _, err := taproot.TweakPublicKey(publicKey, commitment)
	if err != nil {
		t.Errorf("failed to tweak public key: %s", err)
		return
	}

	derivedTweakedPub := ecc.GetPublicKeySchnorr(tweakedPriv)
	if !bytes.Equal(derivedTweakedPub, tweakedPub) {
		t.Errorf(
			"tweaked public key %x is not associated with tweaked private key %x (expected pub %x)",
			tweakedPub, tweakedPriv, derivedTweakedPub,
		)
	}
}
