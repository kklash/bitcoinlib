package signer

import (
	"encoding/hex"
	"testing"

	"github.com/kklash/bitcoinlib/common"
	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/bitcoinlib/ecc"
	"github.com/kklash/bitcoinlib/satutil"
	"github.com/kklash/bitcoinlib/script"
	"github.com/kklash/bitcoinlib/tx"
	"github.com/kklash/bitcoinlib/wif"
)

func mustHex(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return decoded
}

func TestSignP2WPKHTestnet(t *testing.T) {
	constants.CurrentNetwork = constants.BitcoinTestnet
	defer func() {
		constants.CurrentNetwork = constants.BitcoinNetwork
	}()

	privateKey := mustHex("52193c7c8a290a1b93fb140bf3f011ccfc39db77234d9aa7fde059cf011005e9")
	publicKey := ecc.GetPublicKeyCompressed(privateKey)

	var prevOutHash [32]byte
	copy(prevOutHash[:], common.ReverseBytes(mustHex("49e14c722b41b48a1c23d7d306c169cead149475ed21c614a05dd9136e071ab0")))

	outputScript, _ := script.MakeP2WPKHFromPublicKey(publicKey)

	txn := &tx.Tx{
		Version: 2,
		Inputs: []*tx.Input{
			{
				PrevOut:  &tx.PrevOut{Hash: prevOutHash, Index: 1},
				Script:   []byte{},
				Sequence: 0xffffffff,
			},
		},
		Outputs: []*tx.Output{
			{
				Script: outputScript,
				Value:  satutil.BitcoinsToSats(0.0001 - 0.0000013),
			},
		},
	}

	if err := SignInputP2WPKH(txn, 0, privateKey, constants.SigHashAll, satutil.BitcoinsToSats(0.0001)); err != nil {
		t.Errorf("failed to sign transaction: %s", err)
		return
	}

	// self-send testnet transaction: 57f0f7a3eddac5bbfec2fd35b4a9377e5920da18a84d657b1930ad19bc2fd697
	expectedTx := "02000000000101b01a076e13d95da014c621ed759414adce69c106d3d7231c8ab4412b724ce1490100000000ffffffff018e26000000000000160014a859c8253db640d5cdc57a4ee928009ddf2655fa0247304402203aab44ffb427f7c74c2256683959f480af35f5d3ea80189bbe752b09089c2bed022063e0edfc00cfcb578f811e1d21f0e0d5fd05051b9aa2f179abdf7122209db15801210217772ddb0491db40c2cb00cbba1712f8517d5498684fe9ac72b9a199ad8438b900000000"
	actualTx := txn.Hex()

	if actualTx != expectedTx {
		t.Errorf("signed transaction does not match\nWanted %s\nGot    %x", expectedTx, actualTx)
		return
	}
}

func TestSignP2WPKH(t *testing.T) {
	var prevHash [32]byte
	copy(prevHash[:], mustHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))

	var destHash [20]byte
	copy(destHash[:], mustHex("aa4d7985c57e011a8b3dd8e0e5a73aaef41629c5"))

	privateKey, _, _, _ := wif.Decode("KwDiBf89QgGbjEhKnhXJuH7LrciVrZi3qYjgd9M7rFU73sVHnoWn")

	txn := &tx.Tx{
		Version: 1,
		Inputs: []*tx.Input{
			{
				PrevOut:  &tx.PrevOut{Hash: prevHash, Index: 0},
				Script:   []byte{},
				Sequence: 0xffffffff,
			},
		},
		Outputs: []*tx.Output{
			{
				Script: script.MakeP2WPKHFromHash(destHash),
				Value:  10000,
			},
		},
	}

	if err := SignInputP2WPKH(txn, 0, privateKey, constants.SigHashAll, 10000); err != nil {
		t.Errorf("failed to sign transaction: %s", err)
		return
	}

	// This test fixture pulled from bitcoinjs-lib
	expectedTx := "01000000000101ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000ffffffff011027000000000000160014aa4d7985c57e011a8b3dd8e0e5a73aaef41629c502483045022100a8fc5e4c6d7073474eff2af5d756966e75be0cdfbba299518526080ce8b584be02200f26d41082764df89e3c815b8eaf51034a3b68a25f1be51208f54222c1bb6c1601210279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179800000000"
	actualTx := txn.Hex()

	if actualTx != expectedTx {
		t.Errorf("signed transaction does not match\nWanted %s\nGot    %x", expectedTx, actualTx)
		return
	}
}
