package base58check

import (
	"bytes"
	"errors"

	"github.com/kklash/bitcoinlib/base58"
	"github.com/kklash/bitcoinlib/bhash"
)

var (
	// ErrBadBase58Checksum indicates a checksum-mismatch when decoding a base58-check string.
	ErrBadBase58Checksum = errors.New("base58-check decoded checksum does not match")

	// ErrInvalidBase58CheckString is returned by decode when attempting to decode a string which
	// is too short to be a valid base58-check encoded string.
	ErrInvalidBase58CheckStringLength = errors.New("invalid base58-check string length, cannot decode")
)

// Decode decodes the given base58-check string and returns the payload, including potential
// version byte(s). ErrBadBase58Checksum is returned if the checksum does not match.
// ErrInvalidBase58CheckStringLength is returned if the string is less than 4 bytes long, the
// minimum length for a base58-check string.
func Decode(bs58string string) ([]byte, error) {
	bs58decoded, err := base58.Decode(bs58string)
	if err != nil {
		return nil, err
	}

	// 4 byte hash minimum
	if len(bs58decoded) < 4 {
		return nil, ErrInvalidBase58CheckStringLength
	}

	hashed := bhash.DoubleSha256(bs58decoded[:len(bs58decoded)-4])
	if !bytes.Equal(hashed[:4], bs58decoded[len(bs58decoded)-4:]) {
		return nil, ErrBadBase58Checksum
	}

	return bs58decoded[:len(bs58decoded)-4], nil
}
