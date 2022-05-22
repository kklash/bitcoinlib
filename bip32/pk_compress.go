package bip32

import (
	"math/big"

	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/bitcoinlib/ecc"
	"github.com/kklash/ekliptic"
)

// IsCompressedPublicKey returns true if the given byte slice is a compressed public key.
func IsCompressedPublicKey(key []byte) bool {
	return key != nil &&
		len(key) == constants.PublicKeyCompressedLength &&
		(key[0] == constants.PublicKeyCompressedEvenByte ||
			key[0] == constants.PublicKeyCompressedOddByte)
}

// CompressPublicKeyBytes compresses a given uncompressed public key.
// Returns ErrInvalidPublicKey if the key is not of a valid format, or if
// the point it describes is not on the secp256k1 curve.
func CompressPublicKeyBytes(uncompressedPublicKey []byte) ([]byte, error) {
	if err := ValidatePublicKeyBytes(uncompressedPublicKey); err != nil {
		return nil, err
	} else if IsCompressedPublicKey(uncompressedPublicKey) {
		return uncompressedPublicKey, nil
	}

	xBytes, yBytes := uncompressedPublicKey[:32], uncompressedPublicKey[32:]

	x, y := parse256(xBytes), parse256(yBytes)
	return ecc.SerializePointCompressed(x, y), nil
}

// CompressPublicKey compresses a given public key coordinate pair. Returns
// ErrInvalidPublicKey if the point (x, y) is not on the secp256k1 curve.
func CompressPublicKey(x, y *big.Int) ([]byte, error) {
	if err := ValidatePublicKey(x, y); err != nil {
		return nil, err
	}

	return ecc.SerializePointCompressed(x, y), nil
}

// UncompressPublicKey uncompresses a compressed public key and returns
// the (x, y) coordinate pair. Returns ErrInvalidPublicKey if the key is
// of an invalid format, or not on the secp256k1 curve.
func UncompressPublicKey(key []byte) (*big.Int, *big.Int, error) {
	// ValidatePublicKeyBytes needs uncompressPublicKey's logic to determine if a key is valid,
	// so the business logic is separated to prevent infinite recursion.
	if err := ValidatePublicKeyBytes(key); err != nil {
		return nil, nil, err
	}

	x, y := uncompressPublicKey(key)
	return x, y, nil
}

func uncompressPublicKey(key []byte) (*big.Int, *big.Int) {
	x := big.NewInt(0)
	x.SetBytes(key[1:])
	evenY, oddY := ekliptic.Weierstrass(x)

	var y *big.Int
	if key[0]%2 == 0 {
		y = evenY
	} else {
		y = oddY
	}
	return x, y
}
