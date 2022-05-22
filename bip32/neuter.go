package bip32

import "github.com/kklash/bitcoinlib/ecc"

// NeuterCompressed returns the compressed public key of a given private key.
func NeuterCompressed(privateKey []byte) []byte {
	if privateKey == nil {
		return nil
	}

	pubX, pubY := curve.ScalarBaseMult(privateKey)
	return ecc.SerializePointCompressed(pubX, pubY)
}

// NeuterUncompressed returns the uncompressed public key of a given private key.
func NeuterUncompressed(privateKey []byte) []byte {
	if privateKey == nil {
		return nil
	}

	pubX, pubY := curve.ScalarBaseMult(privateKey)
	return ecc.SerializePoint(pubX, pubY)
}

// Neuter returns the compressed or uncompressed public key of
// a given private key, depending on a compression parameter.
func Neuter(privateKey []byte, compressed bool) []byte {
	if compressed {
		return NeuterCompressed(privateKey)
	}

	return NeuterUncompressed(privateKey)
}
