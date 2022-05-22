// package der provides DER encoding for Bitcoin signatures as specified in BIP66.
package der

const (
	// TypeCompound is the prefix byte for any DER encoded signature.
	TypeCompound byte = 0x30

	// TagInteger is the DER tag byte denoting a large integer which follows.
	TagInteger byte = 2

	// MaxIntegerSize is the maximum size of an integer that BIP66 strict-DER encoding allows.
	MaxIntegerSize = 32

	// MinimumSignatureLength and MaximumSignatureLength describe the min/max lengths possible
	// for DER-encoded signatures allowed by BIP66.
	MinimumSignatureLength = 9
	MaximumSignatureLength = 73
)
