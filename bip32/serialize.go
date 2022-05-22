package bip32

import (
	"bytes"
	"errors"

	"github.com/kklash/bitcoinlib/base58check"
	"github.com/kklash/bitcoinlib/constants"
)

// ErrInvalidExtendedKey is returned by Deserialize if it
// is passed a base58 key which fails to decode correctly.
var ErrInvalidExtendedKey = errors.New("Failed to deserialize extended key - Invalid format")

func serialize(key, chainCode, parentFingerprint []byte, depth byte, index, version uint32, isPrivate bool) string {
	buf := new(bytes.Buffer)
	if depth == 0 {
		index = 0
		parentFingerprint = serialize32(0)
	}

	buf.Write(serialize32(version))
	buf.WriteByte(depth)
	buf.Write(parentFingerprint)
	buf.Write(serialize32(index))
	buf.Write(chainCode)

	if isPrivate {
		buf.WriteByte(0)
	}

	buf.Write(key)

	return base58check.Encode(buf.Bytes())
}

// SerializePublic serializes an extended public key to a base58-check encoded string.
func SerializePublic(publicKey, chainCode, parentFingerprint []byte, depth byte, index, version uint32) string {
	return serialize(publicKey, chainCode, parentFingerprint, depth, index, version, false)
}

// SerializePrivate serializes an extended private key to base58-check encoded string.
func SerializePrivate(privateKey, chainCode, parentFingerprint []byte, depth byte, index, version uint32) string {
	return serialize(privateKey, chainCode, parentFingerprint, depth, index, version, true)
}

// Deserialize parses a base58-check encoded extended key (public or private) and returns
// the key, chain code, parent fingerprint, depth from the master key, child index, version
// number prefix, and any error encountered during deserialization. Returns ErrInvalidExtendedKey
// if the key is not of a valid format. Returns ErrInvalidPublicKey if decoding a public key
// and the public key is not on the secp256k1 curve.
func Deserialize(bs58Key string) (key, chainCode, parentFingerprint []byte, depth byte, index, version uint32, err error) {
	var serialized []byte
	serialized, err = base58check.Decode(bs58Key)
	if err != nil {
		err = ErrInvalidExtendedKey
		return
	}

	if len(serialized) != constants.SerializedExtendedKeyLength {
		err = ErrInvalidExtendedKey
		return
	}

	version = parse32(serialized[:4])
	depth = serialized[4]
	parentFingerprint = serialized[5:9]
	index = parse32(serialized[9:13])
	chainCode = serialized[13:45]
	key = serialized[45:]
	if key[0] == 0 {
		// private key, truncate to 32 bytes
		key = key[1:]
	} else {
		if err = ValidatePublicKeyBytes(key); err != nil {
			return
		}
	}

	return
}
