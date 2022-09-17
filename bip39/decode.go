package bip39

import (
	"crypto/sha256"
	"errors"
	"math/big"
)

var (
	// ErrInvalidWordsLength is returned by DecodeWords if the given mnemonic is of
	// an incorrect length.
	ErrInvalidWordsLength = errors.New(
		"invalid BIP39 mnemonic length; must be 12, 15, 18, 21, or 24 words",
	)

	// ErrInvalidWord is returned by DecodeWords if the given mnemonic contains an unknown word.
	ErrInvalidWord = errors.New("BIP39 mnemonic contains word not found in english word list")

	// ErrInvalidChecksum is returned by DecodeWords if the given mnemonic fails checksum validation.
	ErrInvalidChecksum = errors.New("BIP39 mnemonic has invalid checksum")
)

// DecodeWords decodes the given mnemonic into the entropy it encodes, while
// also verifying the length of the mnemonic and its checksum.
//
// Returns any one of ErrInvalidWordsLength, ErrInvalidWord, or ErrInvalidChecksum
// if the mnemonic is not valid.
func DecodeWords(words []string) (entropy []byte, err error) {
	nWords := len(words)
	switch nWords {
	case 12, 15, 18, 21, 24:
	default:
		err = ErrInvalidWordsLength
		return
	}

	payload := new(big.Int)
	for _, word := range words {
		index, ok := WordMap[word]
		if !ok {
			err = ErrInvalidWord
			return
		}

		payload.Lsh(payload, 11)
		payload.Or(payload, big.NewInt(int64(index)))
	}

	nChecksumBits := nWords / 3
	nEntropyBits := nWords*11 - nChecksumBits

	checksumMask := big.NewInt(0xff >> (8 - nChecksumBits))
	checksum := byte(new(big.Int).And(payload, checksumMask).Uint64())

	payload.Rsh(payload, uint(nChecksumBits))
	entropy = payload.FillBytes(make([]byte, nEntropyBits/8))

	hashed := sha256.Sum256(entropy)
	expectedChecksumBits := hashed[0] >> (8 - nChecksumBits)
	if checksum != expectedChecksumBits {
		err = ErrInvalidChecksum
		return
	}

	return
}
