package bech32

import (
	"errors"
	"fmt"

	"github.com/kklash/bits"
)

var (
	// ErrInvalidPayload is returned by Encode if the payload passed is nil or zero-length.
	ErrInvalidPayload = errors.New("Invalid payload for bech32 encoding")

	// ErrInvalidVersion is returned by Encode if the version number is too high to encode as base-32.
	ErrInvalidVersion = fmt.Errorf("Cannot encode version byte higher than %d", len(Alphabet)-1)
)

// encodeValues converts an hrp and a slice of alphabet indeces to a bech32 string.
func encodeValues(hrp string, values []uint5) string {
	values = append(values, bech32CreateChecksum(hrp, values)...)

	bech32 := hrp + Separator
	for i := 0; i < len(values); i++ {
		if int(values[i]) >= len(Alphabet) {
			panic(fmt.Sprintf("cannot bech32 encode value greater than %d", len(Alphabet)-1))
		}

		bech32 += string(Alphabet[values[i]])
	}

	return bech32
}

// bytesToIndeces converts a byte slice into a slice of alphabet indeces by
// splitting the bits of data up into groups of 5.
func bytesToIndeces(data []byte) []uint5 {
	dataBits := bits.BytesToBits(data)
	bitGroups := dataBits.PadRight(BitGroupSize).Split(BitGroupSize)

	values := make([]uint5, len(bitGroups))
	for i, group := range bitGroups {
		// n < (2 ** 5)
		n := group.BigInt().Int64()
		values[i] = uint5(n)
	}

	return values
}

// Encode encodes the given data as a bech32 string, concatenating
// the given human readable prefix, separator, version byte, payload data
// and checksum in base32.
func Encode(hrp string, version byte, data []byte) (string, error) {
	var err error
	if data == nil || len(data) == 0 {
		err = ErrInvalidPayload
	} else if int(version) >= len(Alphabet) {
		err = ErrInvalidVersion
	}

	if err != nil {
		return "", err
	}

	values := bytesToIndeces(data)
	values = append([]uint5{uint5(version)}, values...)
	return encodeValues(hrp, values), nil
}
