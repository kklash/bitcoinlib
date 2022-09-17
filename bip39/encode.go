package bip39

import (
	"crypto/sha256"
	"math/big"
)

var elevenMask = big.NewInt(0b11111111111)

// EncodeToWords encodes the given entropy as a BIP39 mnemonic. Checks the entropy
// is of a valid length (either 16, 20, 24, 28, or 32 bytes).
func EncodeToWords(entropy []byte) ([]string, error) {
	nEntropyBits := len(entropy) * 8
	if err := ValidateEntropySize(nEntropyBits); err != nil {
		return nil, err
	}

	hashedEnt := sha256.Sum256(entropy)
	nChecksumBits := nEntropyBits / 32
	checksumBits := hashedEnt[0] >> (8 - nChecksumBits)

	// payload = entropy || checksum
	payload := new(big.Int).SetBytes(entropy)
	payload.Lsh(payload, uint(nChecksumBits))
	payload.Or(payload, big.NewInt(int64(checksumBits)))

	nWords := (nEntropyBits + nChecksumBits) / 11
	words := make([]string, nWords)

	for i := nWords - 1; i >= 0; i-- {
		index := new(big.Int).And(payload, elevenMask).Int64()
		words[i] = WordList[index]
		payload.Rsh(payload, 11)
	}

	return words, nil
}
