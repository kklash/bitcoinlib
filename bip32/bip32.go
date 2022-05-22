// Package bip32 implements Bitcoin BIP32 deterministic wallet key derivation.
package bip32

import (
	"errors"

	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/ekliptic"
)

var curve = new(ekliptic.Curve)

// ErrInvalidSeed is returned by GenerateMasterKey if the given seed is not valid.
var ErrInvalidSeed = errors.New("Invalid seed; bit size must be between 128 - 512 bits, and a multiple of 32")

// GenerateMasterKey generates a master key and chain code from the
// given seed bytes. Returns ErrInvalidSeed if the seed is not valid.
func GenerateMasterKey(seed []byte) (masterKey, chainCode []byte, err error) {
	if len(seed)%4 != 0 || len(seed) < constants.SeedMinimumSize || len(seed) > constants.SeedMaximumSize {
		err = ErrInvalidSeed
		return
	}

	l := hmacSha512([]byte(constants.BitcoinSeedIV), seed)
	masterKey, chainCode = l[:32], l[32:]
	return
}
