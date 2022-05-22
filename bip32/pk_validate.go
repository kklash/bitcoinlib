package bip32

import (
	"fmt"
	"math/big"

	"github.com/kklash/bitcoinlib/constants"
)

// ErrInvalidPublicKey is returned if a public key is invalid.
// This can mean either that the key is of an invalid length, that the key is nil, or
// that the point described by the key is not on the secp256k1 curve.
var ErrInvalidPublicKey = fmt.Errorf(
	"Invalid public key; must be %d (compressed) or %d (uncompressed) bytes long and must be on secp256k1 curve",
	constants.PublicKeyCompressedLength,
	constants.PublicKeyUncompressedLength,
)

// ValidatePublicKeyBytes checks the validity of a public key, which can be either
// compressed or uncompressed. Returns ErrInvalidPublicKey if the key is of an invalid
// format, or if the point described by the key is not on the secp256k1 curve.
func ValidatePublicKeyBytes(publicKey []byte) error {
	x, y := new(big.Int), new(big.Int)

	if publicKey == nil {
		return ErrInvalidPublicKey
	}

	if IsCompressedPublicKey(publicKey) {
		x, y = uncompressPublicKey(publicKey)
	} else if len(publicKey) == constants.PublicKeyUncompressedLength {
		x.SetBytes(publicKey[:32])
		y.SetBytes(publicKey[32:])
	} else {
		return ErrInvalidPublicKey
	}

	err := ValidatePublicKey(x, y)
	if err != nil {
		return err
	}

	return nil
}

// ValidatePublicKey checks whether the point (x, y) is a valid public key. Returns
// ErrInvalidPublicKey if the key is of an invalid format, or if the point described
// by the key is not on the secp256k1 curve.
func ValidatePublicKey(x, y *big.Int) error {
	if x == nil || y == nil || !curve.IsOnCurve(x, y) {
		return ErrInvalidPublicKey
	}

	return nil
}
