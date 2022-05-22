package script

import (
	"bytes"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/constants"
)

// MakeP2WSHFromHash creates a P2WSH output script using
// the given script hash (hashed once with SHA256).
//  00 <sha256_scripthash>
func MakeP2WSHFromHash(hash [32]byte) []byte {
	script := new(bytes.Buffer)
	script.WriteByte(constants.OP_0)
	script.Write(PushData(hash[:]))
	return script.Bytes()
}

// MakeP2WSHFromScript creates a P2WSH output script using
// the given script, which is hashed once with SHA256 to produce
// the witness program.
func MakeP2WSHFromScript(script []byte) []byte {
	return MakeP2WSHFromHash(bhash.Sha256(script))
}

// IsP2WSH returns whether a byte slice is a valid P2WSH output script.
//  script, _ := hex.DecodeString("0020a745dc51b20631cd78971b67274bc1b1cb06fb800a408f4184fb2bf00dfe9b9d")
//  IsP2WSH(script) // true
//  IsP2WSH(script[1:]) // false
func IsP2WSH(script []byte) bool {
	return script != nil &&
		len(script) == 34 &&
		script[0] == 0x00 &&
		script[1] == 0x20
}

// DecodeP2WSH attempts to decode the given byte slice as a P2WSH script
// pub key. It returns the script-hash contained in the script. Returns
// ErrInvalidScript if the script is not P2WH.
func DecodeP2WSH(script []byte) (hash [32]byte, err error) {
	if !IsP2WSH(script) {
		err = ErrInvalidScript
		return
	}

	copy(hash[:], script[2:34])
	return
}

// WitnessP2WSH builds a P2WSH witness for the given script pub key and redemption
// script (separated into push-only chunks). This witness can be added to a transaction
// as a way of signing the input.
func WitnessP2WSH(scriptPubKey, redeem []byte) ([][]byte, error) {
	chunks, err := Stackify(redeem)
	if err != nil {
		return nil, err
	}

	witnessScriptSig := append(chunks, scriptPubKey)

	return witnessScriptSig, nil
}
