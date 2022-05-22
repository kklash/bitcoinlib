package bip32

import (
	"github.com/kklash/bitcoinlib/bhash"
)

// KeyFingerprint returns the first 4 bytes of the Hash160 of the given public key.
// This is used as the fingerprint for serialized extended keys. Note that publicKey
// can be in compressed or uncompressed form, but will be converted to compressed
// form before hashing. Returns ErrInvalidPublicKey if the key is invalid.
func KeyFingerprint(publicKey []byte) ([]byte, error) {
	if err := ValidatePublicKeyBytes(publicKey); err != nil {
		return nil, ErrInvalidPublicKey
	}

	if !IsCompressedPublicKey(publicKey) {
		var err error
		if publicKey, err = CompressPublicKeyBytes(publicKey); err != nil {
			return nil, err
		}
	}

	h := bhash.Hash160(publicKey) // must be compressed form
	return h[:4], nil
}
