package address

import (
	"github.com/kklash/bitcoinlib/base58check"
	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/constants"
)

// MakeP2PKHFromPublicKey creates a canonical P2PKH address using
// the given compressed or uncompressed public key.
func MakeP2PKHFromPublicKey(publicKey []byte) (string, error) {
	if len(publicKey) != PublicKeyCompressedLength &&
		len(publicKey) != PublicKeyUncompressedLength {
		return "", ErrInvalidPublicKeyLength
	}

	hash := bhash.Hash160(publicKey)
	return MakeP2PKHFromHash(hash), nil
}

// MakeP2PKHFromHash creates a canonical P2PKH address
// using the given public key hash.
func MakeP2PKHFromHash(pkHash [20]byte) string {
	return base58check.EncodeVersion(pkHash[:], constants.CurrentNetwork.PubkeyHash)
}
