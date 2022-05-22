package bip39

import (
	"errors"
	"strings"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bits"
)

var (
	// ErrInvalidMnemonic is returned by Decode if the mnemonic has an invalid length,
	// or if any of its words are not in the English BIP39 word list.
	ErrInvalidMnemonic = errors.New("Unable to decode invalid mnemonic")

	// ErrBadChecksum is returned by Decode if the SHA-256 checksum of the mnemonic does not
	// match the checksum bits appended to the end of the mnemonic.
	ErrBadChecksum = errors.New("Mnemonic checksum is invalid")
)

// Decode decodes the given English BIP39 mnemonic and returns the entropy encoded therein.
// Returns ErrInvalidMnemonic if the mnemonic is of an invalid length, or if any of its words
// do not exist in the word list. Returns ErrBadChecksum if the checksum at the end of the
// entropy does not validate correctly.
func Decode(mnemonic string) ([]byte, error) {
	words := strings.Split(mnemonic, " ")

	if len(words)%3 != 0 {
		return nil, ErrInvalidMnemonic
	}

	bitGroups := make([]bits.Bits, len(words))
	for i, word := range words {
		wordIndex, ok := WordMap[word]
		if !ok {
			return nil, ErrInvalidMnemonic
		}

		bitGroups[i] = bits.UintToBits(uint16(wordIndex)).Trim().PadLeft(11)
	}

	entropyBitsWithChecksum := bits.Join(bitGroups)
	cs := len(entropyBitsWithChecksum) % 32
	csPos := len(entropyBitsWithChecksum) - cs

	entropyBits := entropyBitsWithChecksum[:csPos]
	checksumBits := entropyBitsWithChecksum[csPos:]

	entropy := entropyBits.Bytes()
	hashed := bhash.Sha256(entropy)
	hashBits := bits.BytesToBits(hashed[:cs])

	for i := 0; i < cs; i++ {
		if checksumBits[i] != hashBits[i] {
			return nil, ErrBadChecksum
		}
	}

	return entropy, nil
}
