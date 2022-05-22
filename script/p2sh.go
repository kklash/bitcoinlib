package script

import (
	"bytes"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/constants"
)

// MakeP2SHFromHash creates a P2SH output using
// the given hash of a redeem script.
//  OP_HASH160 <script_hash> OP_EQUAL
func MakeP2SHFromHash(hash [20]byte) []byte {
	scriptPubKey := new(bytes.Buffer)
	scriptPubKey.WriteByte(constants.OP_HASH160)
	scriptPubKey.Write(PushData(hash[:]))
	scriptPubKey.WriteByte(constants.OP_EQUAL)
	return scriptPubKey.Bytes()
}

// MakeP2SHFromScript creates a P2SH output using
// the given script pub key. Makes no attempt to
// validate that the given script is valid.
func MakeP2SHFromScript(script []byte) []byte {
	return MakeP2SHFromHash(bhash.Hash160(script))
}

// IsP2SH returns whether a byte slice is a valid P2SH output script.
//  script, _ := hex.DecodeString("a914aefdfa5219b0ecda7f0f75a2221111769dc848c187")
//  IsP2SH(script) // true
//  IsP2SH(script[2:]) // false
func IsP2SH(scriptPubKey []byte) bool {
	return scriptPubKey != nil &&
		len(scriptPubKey) == 23 &&
		scriptPubKey[0] == constants.OP_HASH160 &&
		scriptPubKey[1] == 0x14 &&
		scriptPubKey[22] == constants.OP_EQUAL
}

// DecodeP2SH attempts to decode the given byte slice as a P2SH script
// pub key. It returns the script-hash contained in the script. Returns
// ErrInvalidScript if the script is not P2SH.
func DecodeP2SH(scriptPubKey []byte) (hash [20]byte, err error) {
	if !IsP2SH(scriptPubKey) {
		err = ErrInvalidScript
		return
	}

	copy(hash[:], scriptPubKey[2:22])
	return
}

// RedeemP2SH builds a P2SH redemption input script using the given
// bitcoin script pub key and the underlying redemption script.
func RedeemP2SH(scriptPubKey, redeem []byte) []byte {
	scriptSig := new(bytes.Buffer)
	scriptSig.Write(redeem)
	scriptSig.Write(PushData(scriptPubKey))
	return scriptSig.Bytes()
}
