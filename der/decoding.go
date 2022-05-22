package der

import (
	"errors"
	"fmt"
	"math/big"
)

var (
	// ErrInvalidSignatureEncoding is returned when decoding a signature fails
	// due to a BIP66 validation problem.
	ErrInvalidSignatureEncoding = errors.New("failed to decode DER signature")
)

func invalidEncodingError(message string, format ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrInvalidSignatureEncoding, fmt.Sprintf(message, format...))
}

func mostSignificantBitFlipped(b byte) bool {
	return b&0x80 == 0x80
}

func hasExtraNullBytes(p []byte) bool {
	return len(p) > 1 && p[0] == 0 && !mostSignificantBitFlipped(p[1])
}

// DecodeSignature decodes the given DER-encoded signature byte slice and returns the
// signature (r, s), as well as the signature hash type. Validates that the format
// matches BIP66 strict-DER encoding and returns ErrInvalidSignatureEncoding otherwise.
func DecodeSignature(derEncoded []byte) (r, s *big.Int, sigHashType uint32, err error) {
	encodedSize := len(derEncoded)

	if encodedSize < MinimumSignatureLength {
		err = invalidEncodingError("signature too small")
		return
	} else if encodedSize > MaximumSignatureLength {
		err = invalidEncodingError("signature too large")
		return
	}

	header := derEncoded[0]
	derMessageSize := int(derEncoded[1])

	if header != TypeCompound {
		err = invalidEncodingError("incorrect DER header byte 0x%x", derEncoded[0])
		return
	} else if derMessageSize != encodedSize-3 { // -3 for TypeCompound, length, and sighash bytes
		err = invalidEncodingError("DER message length doesn't match expected")
		return
	}

	// Decode and validate the encoded signature r value
	rPos := 2
	rTag := derEncoded[rPos]
	rSize := int(derEncoded[rPos+1])

	if rTag != TagInteger {
		err = invalidEncodingError("found incorrect type byte for signature r value")
		return
	} else if rSize > encodedSize-5 || rSize == 0 {
		err = invalidEncodingError("length of r is not valid")
		return
	}

	rBytes := derEncoded[rPos+2 : rPos+2+rSize]

	if mostSignificantBitFlipped(rBytes[0]) {
		err = invalidEncodingError("r cannot be a negative number")
		return
	} else if hasExtraNullBytes(rBytes) {
		err = invalidEncodingError("unnecessary null bytes at start of r")
		return
	}

	// Decode and validate the encoded signature s value
	sPos := 4 + rSize
	sTag := derEncoded[sPos]
	sSize := int(derEncoded[sPos+1])

	// Must have 0x02 integer type
	if sTag != TagInteger {
		err = invalidEncodingError("found incorrect type byte for signature s value")
		return
	} else if sSize+rSize+7 != encodedSize {
		err = invalidEncodingError("sSize + rSize does not match total signature length")
		return
	} else if sSize == 0 {
		err = invalidEncodingError("length of s cannot be zero")
		return
	}

	sBytes := derEncoded[sPos+2 : sPos+2+sSize]
	if mostSignificantBitFlipped(sBytes[0]) {
		err = invalidEncodingError("s cannot be a negative number")
		return
	} else if hasExtraNullBytes(sBytes) {
		err = invalidEncodingError("unnecessary null bytes at start of s")
		return
	}

	r = new(big.Int).SetBytes(rBytes)
	s = new(big.Int).SetBytes(sBytes)
	sigHashType = uint32(derEncoded[sPos+2+sSize])

	return
}
