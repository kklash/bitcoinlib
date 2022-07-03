package bip32

import (
	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/ecc"
)

// KeyFingerprint returns the first 4 bytes of the Hash160 of the given public key.
// This is used as the fingerprint for serialized extended keys. Note that publicKey
// can be in compressed or uncompressed form, but will be converted to compressed
// form before hashing. Returns an error if the public key is not valid.
func KeyFingerprint(publicKey []byte) ([]byte, error) {
	pubX, pubY, err := ecc.DeserializePoint(publicKey)
	if err != nil {
		return nil, err
	}

	if !ecc.IsCompressedPublicKey(publicKey) {
		publicKey = ecc.SerializePointCompressed(pubX, pubY)
	}

	h := bhash.Hash160(publicKey) // must be compressed form
	return h[:4], nil
}
