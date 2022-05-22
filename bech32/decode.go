package bech32

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kklash/bits"
)

var (
	// ErrInvalidBech32 is the base error returned when Decode
	// or Validate is passed a string which is not valid bech32.
	ErrInvalidBech32 = errors.New("Cannot decode invalid bech32 string")

	// ErrInvalidBech32Length wraps ErrInvalidBech32. It is returned
	// when Decode or Validate is called on a string that is of
	// an incorrect length, such that it cannot be valid bech32.
	ErrInvalidBech32Length = fmt.Errorf("%w: invalid length", ErrInvalidBech32)

	// ErrInvalidBech32Character wrapps ErrInvalidBech32. It is returned
	// when Decode or Validate is called on a string containing characters
	// which are not expected to be used in a Bech32 string.
	ErrInvalidBech32Character = fmt.Errorf("%w: invalid character", ErrInvalidBech32)

	// ErrInvalidBech32SeparatorIndex wraps ErrInvalidBech32. It is returned
	// when Decode or Validate is called on a Bech32 string where the
	// separator character is improperly located or missing.
	ErrInvalidBech32SeparatorIndex = fmt.Errorf("%w: separator character not found or in wrong location", ErrInvalidBech32)

	// ErrInvalidBech32MixedCase wrapps ErrInvalidBech32. It is returned
	// when Decode or Validate is called on a string of mixed upper and lower case.
	// Valid bech32 strings must be of either all upper or all lower case.
	ErrInvalidBech32MixedCase = fmt.Errorf("%w: bech32 strings cannot be mixed case", ErrInvalidBech32)

	// ErrInvalidBech32Checksum wrapps ErrInvalidBech32. It is returned
	// when Decode is called on an otherwise valid bech32 string whose
	// checksum fails validation.
	ErrInvalidBech32Checksum = fmt.Errorf("%w: invalid checksum", ErrInvalidBech32)
)

// Validate validates the format of the given bech32 string. The checksum
// is not verified by Validate. Only the formatting of the string is checked.
func Validate(bechAndHrp string) error {
	if len(bechAndHrp) < 8 || len(bechAndHrp) > 90 {
		return ErrInvalidBech32Length
	}

	// The separator must
	// - be present in the string
	// - not be the first char
	// - not be in the last 8 checksum characters
	sepIndex := strings.LastIndex(bechAndHrp, Separator)
	if sepIndex < 1 || sepIndex+ChecksumSize+1 > len(bechAndHrp) {
		return ErrInvalidBech32SeparatorIndex
	}

	for i, c := range []byte(bechAndHrp) {
		if c < 33 || c > 126 {
			// bytes between 33-126 value only
			return ErrInvalidBech32Character
		} else if i > sepIndex && !strings.Contains(Alphabet, string(c)) {
			// Only base32-encoded values after the separator.
			return ErrInvalidBech32Character
		}
	}

	// bech32 strings must be upper or lower case, not mixed
	lowerCase := strings.ToLower(bechAndHrp)
	upperCase := strings.ToUpper(bechAndHrp)
	if bechAndHrp != lowerCase && bechAndHrp != upperCase {
		return ErrInvalidBech32MixedCase
	}

	return nil
}

func separateBechAndHrp(bechAndHrp string) (hrp, bech string) {
	sepIndex := strings.LastIndex(bechAndHrp, Separator)
	hrp = bechAndHrp[:sepIndex]    // Human readable part
	bech = bechAndHrp[sepIndex+1:] // the base32 encoded part (ignore the separator)
	return
}

func bechToBitGroups(hrp, bech string) ([]bits.Bits, error) {
	bitGroups := make([]bits.Bits, len(bech))
	indices := make([]uint5, len(bech))
	for i, c := range []byte(bech) {
		indices[i] = AlphabetIndices[c]
		group := bits.ByteToBits(byte(indices[i]))
		bitGroups[i] = group[8-BitGroupSize:]
	}

	// validate the checksum
	if !bech32VerifyChecksum(hrp, indices) {
		return nil, ErrInvalidBech32Checksum
	}

	return bitGroups, nil
}

// Decode checks the given bech32 string for formatting errors and then attempts to
// decode it. Returns an error wrapping ErrInvalidBech32 if the given string is not
// valid bech32. Returns the human readable prefix, version number, and payload.
func Decode(bechAndHrp string) (hrp string, version byte, data []byte, err error) {
	if err = Validate(bechAndHrp); err != nil {
		return
	}

	bechAndHrp = strings.ToLower(bechAndHrp)

	hrp, bech := separateBechAndHrp(bechAndHrp)
	bitGroups, err := bechToBitGroups(hrp, bech)
	if err != nil {
		return
	}

	if versionBytes := bitGroups[0].BigInt().Bytes(); versionBytes == nil || len(versionBytes) == 0 {
		version = 0
	} else {
		version = versionBytes[0]
	}

	// Cut paylout out from between version byte and checksum
	bitGroups = bitGroups[1 : len(bitGroups)-ChecksumSize]

	// If there is a not-full zero-byte at the end of the bit groups, trim it off
	bitGroups = bits.Join(bitGroups).Split(8)
	if lastGroup := bitGroups[len(bitGroups)-1]; len(lastGroup) != 8 && len(lastGroup.Trim()) == 0 {
		bitGroups = bitGroups[:len(bitGroups)-1]
	}

	bitGroups[len(bitGroups)-1] = bitGroups[len(bitGroups)-1].PadLeft(8)

	allBits := bits.Join(bitGroups)

	data = allBits.Bytes()

	return
}
