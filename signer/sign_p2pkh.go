package wallet

import (
	"github.com/kklash/bitcoinlib/bip32"
	"github.com/kklash/bitcoinlib/script"
	"github.com/kklash/bitcoinlib/tx"
)

func signInputP2PKH(txn *tx.Tx, nInput int, privateKey []byte, sigHashType uint32, compressed bool) error {
	if nInput < 0 || nInput >= len(txn.Inputs) {
		return ErrInputOutOfRange
	}

	publicKey := bip32.Neuter(privateKey, compressed)
	prevOutScript, err := script.MakeP2PKHFromPublicKey(publicKey)
	if err != nil {
		return err
	}

	sigHash, err := txn.SignatureHashForInput(nInput, prevOutScript, sigHashType)
	if err != nil {
		return err
	}

	signature, err := SignSigHash(sigHash[:], privateKey, sigHashType)
	if err != nil {
		return err
	}

	txn.Inputs[nInput].Script = script.RedeemP2PKH(signature, publicKey)
	return nil
}

func SignInputP2PKH(txn *tx.Tx, nInput int, privateKey []byte, sigHashType uint32) error {
	return signInputP2PKH(txn, nInput, privateKey, sigHashType, true)
}

func SignInputP2PKHUncompressed(txn *tx.Tx, nInput int, privateKey []byte, sigHashType uint32) error {
	return signInputP2PKH(txn, nInput, privateKey, sigHashType, false)
}
