package ecc

import (
	"crypto/elliptic"
	"fmt"
	"io"
	"math/big"

	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/ekliptic"
)

// Curve is the secp256k1 curve modeled as an elliptic.Curve interface.
var Curve elliptic.Curve = new(ekliptic.Curve)

// NewPrivateKey generates a private key by reading data from a random source.
// This source should come from a secure high-entropy RNG like crypto/rand.Reader.
func NewPrivateKey(random io.Reader) ([]byte, error) {
	d, err := ekliptic.RandomScalar(random)
	if err != nil {
		return nil, err
	}
	privateKey := d.FillBytes(make([]byte, 32))
	return privateKey, nil
}

// GetPublicKeyCompressed returns the 33-byte compressed public key of a given private key.
func GetPublicKeyCompressed(privateKey []byte) []byte {
	pubX, pubY := Curve.ScalarBaseMult(privateKey)
	return elliptic.MarshalCompressed(Curve, pubX, pubY)
}

// GetPublicKeyUncompressed returns the 65-byte uncompressed public key of a given private key.
func GetPublicKeyUncompressed(privateKey []byte) []byte {
	pubX, pubY := Curve.ScalarBaseMult(privateKey)
	return elliptic.Marshal(Curve, pubX, pubY)
}

// GetPublicKey returns the encoded public key for the given private key, with a boolean parameter
// to decide whether the output public key will be compressed or not.
func GetPublicKey(privateKey []byte, compressed bool) []byte {
	if compressed {
		return GetPublicKeyCompressed(privateKey)
	}
	return GetPublicKeyUncompressed(privateKey)
}

// GetPublicKeySchnorr returns the 32-byte encoded x coordinate of the public key
// belonging to the given private key.
func GetPublicKeySchnorr(privateKey []byte) []byte {
	pubX, _ := Curve.ScalarBaseMult(privateKey)
	return pubX.FillBytes(make([]byte, 32))
}

// CompressPublicKey takes a given public key of any format, deserializes it, and re-encodes
// it in compressed format. Returns ErrPointNotOnCurve if the key is invalid.
func CompressPublicKey(publicKey []byte) ([]byte, error) {
	x, y, err := DeserializePoint(publicKey)
	if err != nil {
		return nil, err
	}
	return SerializePointCompressed(x, y), nil
}

// UncompressPublicKey takes a given public key of any format, deserializes it, and re-encodes
// it in uncompressed format. Returns ErrPointNotOnCurve if the key is invalid.
func UncompressPublicKey(publicKey []byte) ([]byte, error) {
	x, y, err := DeserializePoint(publicKey)
	if err != nil {
		return nil, err
	}
	return elliptic.Marshal(Curve, x, y), nil
}

// IsCompressedPublicKey returns true if the given byte slice appears to be a 33-byte compressed public key.
// It does not check whether the key encodes a valid secp256k1 point.
func IsCompressedPublicKey(key []byte) bool {
	return key != nil &&
		len(key) == constants.PublicKeyCompressedLength &&
		(key[0] == constants.PublicKeyCompressedEvenByte ||
			key[0] == constants.PublicKeyCompressedOddByte)
}

// SumPrivateKeys combines any number of private keys together into one key which
// can form an aggregate signature.
//
// Returns an error if any of the keys are not valid secp256k1 scalar values.
func SumPrivateKeys(privateKeys ...[]byte) ([]byte, error) {
	sum := new(big.Int)
	keyInt := new(big.Int)
	for _, key := range privateKeys {
		keyInt.SetBytes(key)
		if !ekliptic.IsValidScalar(keyInt) {
			return nil, fmt.Errorf("SumPrivateKeys used on invalid private key")
		}
		sum.Add(sum, keyInt)
	}
	sumBytes := sum.FillBytes(make([]byte, 32))
	return sumBytes, nil
}

// SumPublicKeys combines any number of schnorr public keys together into one key which
// can be used to verify signatures aggregated from the component keys.
//
// Returns an error if any of the keys are not schnorr public keys.
func SumPublicKeys(publicKeys ...[]byte) ([]byte, error) {
	sumX := new(big.Int)
	sumY := new(big.Int)
	for _, key := range publicKeys {
		if len(key) != constants.PublicKeySchnorrLength {
			return nil, fmt.Errorf(
				"CombinePublicKeys only works with %d-byte schnorr public keys",
				constants.PublicKeySchnorrLength,
			)
		}

		x, y, err := DeserializePoint(key)
		if err != nil {
			return nil, err
		}

		sumX, sumY = ekliptic.AddAffine(sumX, sumY, x, y)
	}

	sumPub := sumX.FillBytes(make([]byte, 32))
	return sumPub, nil
}
