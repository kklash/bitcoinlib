package address

import (
	"github.com/kklash/bitcoinlib/bech32"
	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/constants"
)

// MakeP2WPKHFromPublicKey creates a P2WPKH segwit address using the given
// compressed public key.
func MakeP2WPKHFromPublicKey(publicKey []byte) (string, error) {
	if len(publicKey) != constants.PublicKeyCompressedLength {
		return "", ErrInvalidPublicKeyLength
	}

	witnessProgram := bhash.Hash160(publicKey)
	address, err := MakeP2WPKHFromHash(witnessProgram)
	if err != nil {
		return "", err
	}

	return address, nil
}

// MakeP2WPKHFromHash creates a P2WPKH segwit address using the given public key hash.
func MakeP2WPKHFromHash(pkHash [20]byte) (string, error) {
	if len(constants.CurrentNetwork.Bech32) == 0 {
		return "", ErrNoSegwitSupport
	}

	address, err := bech32.Encode(
		constants.CurrentNetwork.Bech32,
		constants.WitnessVersionZero,
		pkHash[:],
	)

	if err != nil {
		return "", err
	}

	return address, nil
}
