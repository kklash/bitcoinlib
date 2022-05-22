package der

import (
	"bytes"
	"errors"
	"math/big"
)

var (
	// ErrNotEncodable is returned when an encoding function is passed a value
	// which cannot be encoded according to strict-DER as specified in BIP66.
	ErrNotEncodable = errors.New("failed to encode value in strict-DER")

	// ErrInvalidSigHashType is returned when DER-encoding a signature, if the sighash type
	// given is larger than can fit in a byte.
	ErrInvalidSigHashType = errors.New("failed to encode invalid sighash type in strict-DER")
)

// EncodeBigInt encodes a big.Int as a DER byte slice.
// Returns ErrNotEncodable if it is not encodable.
func EncodeBigInt(v *big.Int) ([]byte, error) {
	if err := CheckEncodableBigInt(v); err != nil {
		return nil, err
	}

	vBytes := v.Bytes()

	// the first bit being flipped means a negative number, which is
	// disallowed. We add an empty padding byte in this case.
	if len(vBytes) == 0 || vBytes[0]&0x80 == 0x80 {
		vBytes = append([]byte{0}, vBytes...)
	}

	encoded := append(
		[]byte{TagInteger, byte(len(vBytes))},
		vBytes...,
	)

	return encoded, nil
}

// CheckEncodableBigInt validates whether the big.Int v can be encoded in strict-DER as defined by BIP66.
// Returns ErrNotEncodable if it is not encodable.
func CheckEncodableBigInt(v *big.Int) error {
	if v == nil || v.BitLen() > MaxIntegerSize*8 || v.Sign() == -1 {
		return ErrNotEncodable
	}

	return nil
}

// EncodeSignature encodes the given signature
// Returns ErrNotEncodable if the signature is not encodable in strict-DER.
// Returns ErrInvalidSigHashType if sigHashType cannot encode as a single byte.
func EncodeSignature(r, s *big.Int, sigHashType uint32) ([]byte, error) {
	if sigHashType > 0xff {
		return nil, ErrInvalidSigHashType
	}

	rEncoded, err := EncodeBigInt(r)
	if err != nil {
		return nil, err
	}

	sEncoded, err := EncodeBigInt(s)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	buf.WriteByte(TypeCompound)
	buf.WriteByte(byte(len(rEncoded) + len(sEncoded)))
	buf.Write(rEncoded)
	buf.Write(sEncoded)
	buf.WriteByte(byte(sigHashType))

	return buf.Bytes(), nil
}
