package script

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/kklash/bitcoinlib/constants"
)

// Uses example from https://github.com/bitcoin-core/btcdeb/blob/1e7ff0f8da533bb8cf8ff6ebb4f160852cbb3685/doc/tapscript-example-with-tap.md
func TestMakeP2TR(t *testing.T) {
	// script1:
	//  144
	//  OP_CHECKSEQUENCEVERIFY
	//  OP_DROP
	//  <9997a497d964fc1a62885b05a51166a65a90df00492c8d7cf61d6accf54803be>
	//  OP_CHECKSIG
	alicePub, _ := hex.DecodeString("9997a497d964fc1a62885b05a51166a65a90df00492c8d7cf61d6accf54803be")
	script1 := append(
		PushNumber(144),
		constants.OP_CHECKSEQUENCEVERIFY,
		constants.OP_DROP,
	)
	script1 = append(
		script1,
		PushData(alicePub)...,
	)
	script1 = append(script1, constants.OP_CHECKSIG)

	// script2:
	//  OP_SHA256
	//  <6c60f404f8167a38fc70eaf8aa17ac351023bef86bcb9d1086a19afe95bd5333>
	//  OP_EQUALVERIFY
	//  <4edfcf9dfe6c0b5c83d1ab3f78d1b39a46ebac6798e08e19761f5ed89ec83c10>
	//  OP_CHECKSIG
	bobPub, _ := hex.DecodeString("4edfcf9dfe6c0b5c83d1ab3f78d1b39a46ebac6798e08e19761f5ed89ec83c10")
	hashed, _ := hex.DecodeString("6c60f404f8167a38fc70eaf8aa17ac351023bef86bcb9d1086a19afe95bd5333")

	script2 := append([]byte{constants.OP_SHA256}, PushData(hashed)...)
	script2 = append(script2, constants.OP_EQUALVERIFY)
	script2 = append(script2, PushData(bobPub)...)
	script2 = append(script2, constants.OP_CHECKSIG)

	mast := &MastBranch{
		Left: &MastLeaf{
			Version: constants.TaprootLeafVersionTapscript,
			Script:  script1,
		},
		Right: &MastLeaf{
			Version: constants.TaprootLeafVersionTapscript,
			Script:  script2,
		},
	}

	mastRootHash := mast.Hash()
	expectedRootHash, _ := hex.DecodeString("41646f8c1fe2a96ddad7f5471bc4fee7da98794ef8c45a4f4fc6a559d60c9f6b")
	if !bytes.Equal(mastRootHash[:], expectedRootHash) {
		t.Errorf("MAST root hash does not match\nWanted %x\nGot    %x", expectedRootHash, mastRootHash)
		return
	}

	internalPublicKey, _ := hex.DecodeString("f30544d6009c8d8d94f5d030b2e844b1a3ca036255161c479db1cca5b374dd1c")
	scriptPubKey, err := MakeP2TR(internalPublicKey, mast)
	if err != nil {
		t.Errorf("failed to build P2TR output script: %s", err)
		return
	}

	expectedScript, _ := hex.DecodeString("5120a5ba0871796eb49fb4caa6bf78e675b9455e2d66e751676420f8381d5dda8951")
	if !bytes.Equal(scriptPubKey, expectedScript) {
		t.Errorf("P2TR script does not match\nWanted %x\nGot    %x", expectedScript, scriptPubKey)
		return
	}
}
