// Package wif provides APIs to encode and decode
// Bitcoin private keys in Wallet Import Format.
package wif

import (
	"errors"

	"github.com/kklash/bitcoinlib/base58check"
)

var (
	// ErrInvalidPrivateKey indicates nil or a byte slice of an incorrect length
	// was passed to Encode or EncodeUncompressed.
	ErrInvalidPrivateKey = errors.New("cannot make WIF key, expected private key of 32 bytes")

	// ErrInvalidWifKey is returned by Decode when a valid base58check
	// string is passed, but the decoded string contents are not WIF.
	ErrInvalidWifKey = errors.New("key is not in wallet import format")
)

func encode(privkey []byte, version byte, compressed bool) (string, error) {
	if privkey == nil || len(privkey) != 32 {
		return "", ErrInvalidPrivateKey
	}

	if compressed {
		privkey = append(privkey, 0x01)
	}

	encoded := base58check.EncodeVersion(privkey, uint16(version))
	return encoded, nil
}

// Encode converts a private key into WIF format.
func Encode(privkey []byte, version byte) (string, error) {
	return encode(privkey, version, true)
}

// EncodeUncompressed converts a private key into WIF format, specifying
// that the key belongs to an uncompressed public key.
func EncodeUncompressed(privkey []byte, version byte) (string, error) {
	return encode(privkey, version, false)
}

// Decode decodes the WIF private key string and returns the
// raw private key, the version byte, and a bool indicating
// if the key corresponds to a compressed public key.
func Decode(wifkey string) ([]byte, byte, bool, error) {
	wifDecoded, err := base58check.Decode(wifkey)
	if err != nil {
		return nil, 0, false, err
	}

	// 32 byte key + 1 byte version (+ 1 byte compressed flag)
	if len(wifDecoded) < 33 || len(wifDecoded) > 34 {
		return nil, 0, false, ErrInvalidWifKey
	}

	version := wifDecoded[0]
	privkey := wifDecoded[1:33]

	var compressed bool = false
	if len(wifDecoded) == 34 {
		if wifDecoded[33] == 0x01 {
			compressed = true
		} else { // unexpected extra byte
			return nil, 0, false, ErrInvalidWifKey
		}
	}

	return privkey, version, compressed, nil
}

// Validate returns true if the given string is a valid WIF key.
func Validate(wifkey string) bool {
	_, _, _, err := Decode(wifkey)
	return err == nil
}
