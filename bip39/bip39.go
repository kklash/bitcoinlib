// Package bip39 implements the Bitcoin BIP39 mnemonic phrase encoding standard.
package bip39

import (
	"crypto/sha512"
	"errors"
	"io"

	"github.com/kklash/bitcoinlib/constants"
	"golang.org/x/crypto/pbkdf2"
)

// The flow goes
// entropy -> mnemonic
// mnemonic + password -> seed
// seed -> bip32 key and chain code

// ErrInvalidEntropySize is returned by ValidateEntropySize and GenerateEntropy if
// the entropy's bit size is not within the bounds of constants.EntropyMinimumSize and constants.EntropyMaximumSize.
var ErrInvalidEntropySize = errors.New("Entropy bit size; must be between 128 - 256 bits and must be divisible by 32")

// ValidateEntropySize checks if entropy of bitSize bits is valid
// and returns ErrInvalidEntropySize if it is not.
func ValidateEntropySize(bitSize int) error {
	if bitSize < constants.EntropyMinimumSize || bitSize > constants.EntropyMaximumSize || bitSize%32 != 0 {
		return ErrInvalidEntropySize
	}

	return nil
}

// GenerateEntropy creates entropy of the given bitSize. Returns ErrInvalidEntropySize if size is not valid.
func GenerateEntropy(rand io.Reader, bitSize int) ([]byte, error) {
	if err := ValidateEntropySize(bitSize); err != nil {
		return nil, err
	}

	entropy := make([]byte, bitSize/8)
	if _, err := io.ReadFull(rand, entropy); err != nil {
		return nil, err
	}

	return entropy, nil
}

// GenerateMnemonic creates a random BIP39 mnemonic phrase using entropy from rand.
func GenerateMnemonic(rand io.Reader, nWords int) (string, error) {
	bitSize := nWords * 32 / 3
	entropy, err := GenerateEntropy(rand, bitSize)
	if err != nil {
		return "", err
	}

	mnemonic, err := Encode(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

// DeriveSeed turns the given mnemonic and passphrase into a seed usable for BIP32 master key derivation.
func DeriveSeed(mnemonic, passphrase string) []byte {
	salt := []byte(constants.SeedSaltPrefix + passphrase)
	return pbkdf2.Key([]byte(mnemonic), salt, constants.SeedIterationCount, constants.SeedSize/8, sha512.New)
}
