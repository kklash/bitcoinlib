package address

import (
	"github.com/kklash/bitcoinlib/base58check"
	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/constants"
)

// MakeP2SHFromScript creates a P2SH address using the given script.
// WARNING No script validation is performed in this
// function. You won't be able to recover any coins to an
// address made with an invalid script.
func MakeP2SHFromScript(script []byte) string {
	scriptHash := bhash.Hash160(script)
	return MakeP2SHFromHash(scriptHash)
}

// MakeP2SHFromHash creates a P2SH address using the
// given script hash.
func MakeP2SHFromHash(scriptHash [20]byte) string {
	return base58check.EncodeVersion(scriptHash[:], constants.CurrentNetwork.ScriptHash)
}
