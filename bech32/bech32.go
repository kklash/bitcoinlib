// Package bech32 provides base-32 encoding as specified in BIP173.
package bech32

import (
	"github.com/kklash/bitcoinlib/constants"
)

const (
	// Alphabet is the base 32 alphabet used for encoding.
	Alphabet = constants.Bech32Alphabet

	// BitGroupSize is the bit size of groups which bytes are split into before encoding into base-32.
	BitGroupSize = 5

	// ChecksumSize is the encoded length of the checksum appended to bech32-encoded strings.
	ChecksumSize = 6

	// Separator is the separating character which joins the HRP with the version number and encoded data.
	Separator = constants.Bech32Separator
)

// AlphabetIndices is a mapping of every char in Alphabet to their index in Alphabet.
var AlphabetIndices = make(map[byte]uint5)

// uint5 is a type declaration to make things more clear when
// we are referring to numbers which should be less than 2^5.
type uint5 uint8

func init() {
	for i := 0; i < len(Alphabet); i++ {
		r := Alphabet[i]
		AlphabetIndices[r] = uint5(i)
	}
}
