package signer

import (
	"github.com/kklash/bitcoinlib/ecc"
	"github.com/kklash/bitcoinlib/script"
	"github.com/kklash/bitcoinlib/tx"
)

func SignInputP2WPKH(txn *tx.Tx, nInput int, privateKey []byte, sigHashType uint32, inputValue uint64) error {
	if nInput < 0 || nInput >= len(txn.Inputs) {
		return ErrInputOutOfRange
	}

	publicKey := ecc.GetPublicKeyCompressed(privateKey)

	// The prevOutScript is not the literal witness program of the prevout.
	// Instead it is the script pub key implied by the witness program, hence
	// we supply a canonical P2PKH script pub key as the prevOutScript.
	prevOutScript, err := script.MakeP2PKHFromPublicKey(publicKey)
	if err != nil {
		return err
	}

	sigHash, err := txn.SignatureHashForWitnessInput(nInput, prevOutScript, sigHashType, inputValue)
	if err != nil {
		return err
	}

	signature, err := SignSigHash(sigHash[:], privateKey, sigHashType)
	if err != nil {
		return err
	}

	if txn.Witnesses == nil {
		txn.Witnesses = make([]tx.Witness, len(txn.Inputs))
	}

	for i := 0; i < len(txn.Inputs); i++ {
		if i == nInput {
			txn.Witnesses[i] = script.WitnessP2WPKH(signature, publicKey)
		} else if txn.Witnesses[i] == nil {
			txn.Witnesses[i] = tx.Witness{}
		}
	}

	// Segwit signatures use empty input scripts
	txn.Inputs[nInput].Script = []byte{}

	return nil
}
