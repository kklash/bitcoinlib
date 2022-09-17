// Package bip39 implements the BIP39 mnemonic encoding scheme for private key seed data.
//
// To generate a mnemonic, generate some random entropy (either 16, 20, 24, 28, or 32 bytes),
// and use EncodeToWords to encode the entropy as a mnemonic. Normally this phrase is given
// to the user as a backup.
//
// To utilize the mnemonic to generate a seed for cryptographic keys, pass the mnemonic and
// an optional deniability passphrase to DeriveSeed.
//
//	https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
package bip39

import (
	"errors"
	"io"
)

// The flow goes
// entropy -> mnemonic
// mnemonic + password -> seed
// seed -> bip32 key and chain code

// ErrInvalidEntropySize is returned by ValidateEntropySize and GenerateEntropy if
// the entropy's bit size is not within the bounds of EntropyMinimumSize and EntropyMaximumSize.
var ErrInvalidEntropySize = errors.New(
	"bad entropy bit size; must be between 128 - 256 bits and must be divisible by 32",
)

const (
	// EntropyMinimumSize and EntropyMaximumSize specify the minimum/maximum number of bits of
	// entropy needed to generate a secure mnemonic.
	EntropyMinimumSize = 128
	EntropyMaximumSize = 256
)

// ValidateEntropySize checks if entropy of bitSize bits is valid
// and returns ErrInvalidEntropySize if it is not.
func ValidateEntropySize(bitSize int) error {
	if bitSize < EntropyMinimumSize || bitSize > EntropyMaximumSize || bitSize%32 != 0 {
		return ErrInvalidEntropySize
	}

	return nil
}

// GenerateEntropy creates entropy of the given bitSize. Returns
// ErrInvalidEntropySize if size is not valid.
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
func GenerateMnemonic(rand io.Reader, nWords int) ([]string, error) {
	bitSize := nWords * 32 / 3
	entropy, err := GenerateEntropy(rand, bitSize)
	if err != nil {
		return nil, err
	}

	words, err := EncodeToWords(entropy)
	if err != nil {
		return nil, err
	}

	return words, nil
}
