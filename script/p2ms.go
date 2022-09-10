package script

import (
	"bytes"

	"github.com/kklash/bitcoinlib/constants"
)

// MakeP2MS creates a pay-to-M-of-N-multisig output script. The
// sigsRequired parameter specifies M, and you can supply as many
// public keys as you want as long as the number of public keys >= M.
//
//	M <pubkey_1> <pubkey_2> ... <pubkey_N> N OP_CHECKMULTISIG
func MakeP2MS(sigsRequired uint32, publicKeys ...[]byte) []byte {
	if int(sigsRequired) > len(publicKeys) {
		panic("cannot create P2MS script requiring more signatures than the number of available public keys")
	} else if sigsRequired == 0 {
		panic("sigsRequired for P2MS script cannot be zero")
	}

	scriptPubKey := new(bytes.Buffer)
	scriptPubKey.Write(PushNumber(int64(sigsRequired)))
	for _, publicKey := range publicKeys {
		scriptPubKey.Write(PushData(publicKey))
	}
	scriptPubKey.Write(PushNumber(int64(len(publicKeys))))
	scriptPubKey.WriteByte(constants.OP_CHECKMULTISIG)
	return scriptPubKey.Bytes()
}

// RedeemP2MS builds a P2MS redemption script using the given DER-encoded signatures.
func RedeemP2MS(signatures ...[]byte) []byte {
	if len(signatures) == 0 {
		panic("cannot redeem P2MS script with no signatures")
	}

	scriptSig := new(bytes.Buffer)
	scriptSig.WriteByte(constants.OP_0)
	for _, signature := range signatures {
		scriptSig.Write(PushData(signature))
	}
	return scriptSig.Bytes()
}
