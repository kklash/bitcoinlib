// Package constants provides magic numbers and other useful data, such as opcodes,
// which would otherwise have to be copy/pasted and re-declared across packages.
package constants

const (
	// BitcoinSeedIV is the initialization vector used for creating bitcoin wallet seeds.
	BitcoinSeedIV = "Bitcoin seed"

	// Base58Alphabet is the Bitcoin Base 58 alphabet.
	Base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	// Bech32Alphabet is the base-32 encoding alphabet for BECH32 as specified in BIP-173.
	Bech32Alphabet = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

	// Bech32Separator is the separating character in bech32 which separates the HRP from the version number and encoded data.
	Bech32Separator = "1"

	// Bip32Hardened is the HD key index threshold above which any derived child keys are hardened.
	Bip32Hardened uint32 = 0x80000000

	// BlockMaxSize is the maximum block size, not including segwit data.
	BlockMaxSize = 1000000

	// OpReturnMaxSize is the maximum size of an OP_RETURN output data payload.
	OpReturnMaxSize = 80

	// SatoshisPerBitcoin is the number of base satoshi units per BTC
	SatoshisPerBitcoin = 100_000_000

	// SeedMinimumSize and SeedMaximumSize are the lower and upper limits on
	// the byte-size of a BIP39 seed which can be used for generating a master key.
	SeedMinimumSize int = 128 / 8
	SeedMaximumSize int = 512 / 8

	// SerializedExtendedKeyLength is the byte length of a bip32 serialized extended key.
	SerializedExtendedKeyLength int = 78

	// SigHash enum types
	SigHashAll          uint32 = 1
	SigHashNone         uint32 = 2
	SigHashSingle       uint32 = 3
	SigHashAnyoneCanPay uint32 = 0x80

	// PublicKeyCompressedLength is the byte-size of a compressed public key.
	PublicKeyCompressedLength int = 33
	// PublicKeyUncompressedLength is the byte-size of an uncompressed public key.
	PublicKeyUncompressedLength int = 65
	// PublicKeySchnorrLength is the byte-size of a BIP340 schnorr public key.
	PublicKeySchnorrLength int = 32

	// PublicKeyCompressedEvenByte and PublicKeyCompressedOddByte indicate the bytes prefixing
	// compressed public keys to indicate whether the Y coordinate of the public key point (X, Y) is even or odd.
	PublicKeyCompressedEvenByte byte = 2
	PublicKeyCompressedOddByte  byte = 3
	PublicKeyUncompressedPrefix byte = 4

	// TaprootLeafVersionTapscript is the version number used in Taproot Tapscript leaf nodes.
	TaprootLeafVersionTapscript = 0xc0

	// WitnessVersionZero is the first witness version introduced. It is used for bech32-encoded
	// P2WPKH and P2WSH witness programs.
	WitnessVersionZero = 0
)

// AddressFormat is used to describe different standardized script pubkey formats.
type AddressFormat string

const (
	FormatP2PKH       AddressFormat = "P2PKH"
	FormatP2SH        AddressFormat = "P2SH"
	FormatP2WPKH      AddressFormat = "P2WPKH"
	FormatP2WSH       AddressFormat = "P2WSH"
	FormatNONSTANDARD AddressFormat = "NONSTANDARD"
)

var (
	// Bech32ChecksumGen are the seed numbers used to create bech32 checksums.
	// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki#Checksum
	Bech32ChecksumGen = [...]int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

	// TxSegwitFlag is the byte sequence which, if present in a serialized transaction at a certain
	// location, indicates the precense of witness data.
	TxSegwitFlag = [...]byte{0, 1}
)
