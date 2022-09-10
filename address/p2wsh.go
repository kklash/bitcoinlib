package address

import (
	"github.com/kklash/bitcoinlib/bech32"
	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/constants"
)

// MakeP2WSHFromScript creates a P2WSH address out of a given script.
//
// WARNING No script validation is performed in this
// function. You won't be able to recover any coins to an
// adress made with an invalid script.
func MakeP2WSHFromScript(script []byte) (string, error) {
	witnessProgram := bhash.Sha256(script)
	address, err := MakeP2WSHFromHash(witnessProgram)

	if err != nil {
		return "", err
	}

	return address, nil
}

// MakeP2WSHFromHash creates a P2WSH address using a given SHA256 script hash.
func MakeP2WSHFromHash(scriptHash [32]byte) (string, error) {
	if len(constants.CurrentNetwork.Bech32) == 0 {
		return "", ErrNoSegwitSupport
	}

	address, err := bech32.Encode(
		constants.CurrentNetwork.Bech32,
		constants.WitnessVersionZero,
		scriptHash[:],
	)

	if err != nil {
		return "", err
	}

	return address, nil
}
