package script

import (
	"bytes"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/constants"
)

// MakeP2PKHFromHash builds a canonical P2PKH output
// script using the given public key hash.
//
//	OP_DUP OP_HASH160 <pubkey_hash> OP_EQUALVERIFY OP_CHECKSIG
func MakeP2PKHFromHash(hash [20]byte) []byte {
	script := new(bytes.Buffer)
	script.WriteByte(constants.OP_DUP)
	script.WriteByte(constants.OP_HASH160)
	script.Write(PushData(hash[:]))
	script.WriteByte(constants.OP_EQUALVERIFY)
	script.WriteByte(constants.OP_CHECKSIG)
	return script.Bytes()
}

// MakeP2PKHFromPublicKey builds a canonical P2PKH output
// script using the given public key. Returns ErrInvalidPublicKeyLength
// if the public key provided is not valid.
func MakeP2PKHFromPublicKey(publicKey []byte) ([]byte, error) {
	if len(publicKey) != constants.PublicKeyCompressedLength &&
		len(publicKey) != constants.PublicKeyUncompressedLength {
		return nil, ErrInvalidPublicKeyLength
	}

	return MakeP2PKHFromHash(bhash.Hash160(publicKey)), nil
}

// Is P2PKH returns whether the given byte slice is a valid P2PKH output script.
//
//	script, _ := hex.DecodeString("76a914c41c836560406c6169537f9dc9520184879f03e288ac")
//	IsP2PKH(script) // true
//	IsP2PKH(script[1:]) // false
func IsP2PKH(script []byte) bool {
	return script != nil &&
		len(script) == 25 &&
		script[0] == constants.OP_DUP &&
		script[1] == constants.OP_HASH160 &&
		script[2] == 0x14 &&
		script[23] == constants.OP_EQUALVERIFY &&
		script[24] == constants.OP_CHECKSIG
}

// DecodeP2PKH attempts to decode the given byte slice as a P2PKH script
// pub key. It returns the pubkey-hash contained in the script. Returns
// ErrInvalidScript if the script is not P2PKH.
func DecodeP2PKH(script []byte) (hash [20]byte, err error) {
	if !IsP2PKH(script) {
		err = ErrInvalidScript
		return
	}

	copy(hash[:], script[3:23])
	return
}

// RedeemP2PKH builds a canonical P2PKH redemption input script using the given
// DER encoded signature and public key.
func RedeemP2PKH(signature, publicKey []byte) []byte {
	return append(PushData(signature), PushData(publicKey)...)
}
