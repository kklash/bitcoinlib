package signer

import (
	"github.com/kklash/bitcoinlib/der"
	"github.com/kklash/bitcoinlib/ecc"
)

func SignSigHash(hash, privateKey []byte, sigHashType uint32) ([]byte, error) {
	r, s := ecc.SignECDSA(privateKey, hash)

	signature, err := der.EncodeSignature(r, s, sigHashType)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
