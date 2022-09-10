package script

import (
	"bytes"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/constants"
)

// MakeP2WPKHFromHash creates an output script for P2WPKH using the given
// public key hash. Note that you MUST use a compressed public key or
// you will not be able to redeem any UTXOs sent to this output script.
//
//	00 <pubkey_hash>
func MakeP2WPKHFromHash(hash [20]byte) []byte {
	script := new(bytes.Buffer)
	script.WriteByte(constants.OP_0)
	script.Write(PushData(hash[:]))
	return script.Bytes()
}

// MakeP2WPKHFromPublicKey creates an output script for P2WPKH using the given
// public key. Returns ErrInvalidScript if you do not pass a compressed public key.
func MakeP2WPKHFromPublicKey(publicKey []byte) ([]byte, error) {
	if len(publicKey) != constants.PublicKeyCompressedLength {
		return nil, ErrInvalidPublicKeyLength
	}

	return MakeP2WPKHFromHash(bhash.Hash160(publicKey)), nil
}

// IsP2WPKH returns whether a byte slice is a valid P2WPKH output script.
//
//	script, _ := hex.DecodeString("0014b21658ea3960030b4dc4634f51ecae0505daf2e6")
//	IsP2WPKH(script) // true
//	IsP2WPKH(script[1:]) // false
func IsP2WPKH(script []byte) bool {
	return script != nil &&
		len(script) == 22 &&
		script[0] == 0x00 &&
		script[1] == 0x14
}

// DecodeP2WPKH attempts to decode the given byte slice as a P2WPKH script
// pub key. It returns the pubkey-hash contained in the script. Returns
// ErrInvalidScript if the script is not P2WPKH.
func DecodeP2WPKH(script []byte) (hash [20]byte, err error) {
	if !IsP2WPKH(script) {
		err = ErrInvalidScript
		return
	}

	copy(hash[:], script[2:22])
	return
}

// WitnessP2WPKH builds a P2WPKH witness for the given signature and public key.
// This witness can be added to a transaction as a way of signing the input.
func WitnessP2WPKH(signature, publicKey []byte) [][]byte {
	return [][]byte{signature, publicKey}
}
