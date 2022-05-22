package bip39

import (
	"errors"
	"strings"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bits"
)

// ErrInvalidData is returned by EncodeAny if len(data) % 4 != 0.
// This is a mathematical limitation imposed by the BIP39 spec.
var ErrInvalidData = errors.New("Data cannot be encoded if its byte size is not divisible by 4")

// Encode encodes the given entropy as an English BIP39 mnemonic. Returns ErrInvalidEntropySize
// if the size of entropy is not valid.
func Encode(entropy []byte) (string, error) {
	if err := ValidateEntropySize(len(entropy) * 8); err != nil {
		return "", err
	}

	mnemonic, err := EncodeAny(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

// EncodeAny encodes arbitrary data into a BIP39 mnemonic string without imposing
// restrictions on the size of the data.
func EncodeAny(data []byte) (string, error) {
	if len(data)%4 != 0 {
		return "", ErrInvalidData
	}

	hashed := bhash.Sha256(data)
	dataBits := bits.BytesToBits(data)
	cs := len(dataBits) / 32
	hashBits := bits.BytesToBits(hashed[:cs])
	dataBits = append(dataBits, hashBits[:cs]...)

	var bitGroups []bits.Bits = dataBits.Split(11)
	words := make([]string, len(bitGroups))

	for i := 0; i < len(bitGroups); i++ {
		v := bitGroups[i].BigInt()
		words[i] = WordList[v.Uint64()]
	}

	return strings.Join(words, " "), nil
}
