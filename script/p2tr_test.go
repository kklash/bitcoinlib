package script

import (
	"bytes"
	"testing"

	"github.com/kklash/bitcoinlib/constants"
)

func TestMakeP2TR(t *testing.T) {
	type Fixture struct {
		internalPublicKey []byte
		scriptTree        Hasher
		scriptPubKey      []byte
	}

	// From BIP341 test vectors.
	// https://github.com/bitcoin/bips/blob/52f68fecd8ec9604672e26392468e7e7edf25a5e/bip-0341/wallet-test-vectors.json
	fixtures := []*Fixture{
		{
			internalPublicKey: hex2bytes("d6889cb081036e0faefa3a35157ad71086b123b2b144b649798b494c300a961d"),
			scriptTree:        nil,
			scriptPubKey:      hex2bytes("512053a1f6e454df1aa2776a2814a721372d6258050de330b3c6d10ee8f4e0dda343"),
		},

		{
			internalPublicKey: hex2bytes("187791b6f712a8ea41c8ecdd0ee77fab3e85263b37e1ec18a3651926b3a6cf27"),
			scriptTree: &MastLeaf{
				Version: constants.TaprootLeafVersionTapscript,
				Script:  hex2bytes("20d85a959b0290bf19bb89ed43c916be835475d013da4b362117393e25a48229b8ac"),
			},
			scriptPubKey: hex2bytes("5120147c9c57132f6e7ecddba9800bb0c4449251c92a1e60371ee77557b6620f3ea3"),
		},

		{
			internalPublicKey: hex2bytes("93478e9488f956df2396be2ce6c5cced75f900dfa18e7dabd2428aae78451820"),
			scriptTree: &MastLeaf{
				Version: constants.TaprootLeafVersionTapscript,
				Script:  hex2bytes("20b617298552a72ade070667e86ca63b8f5789a9fe8731ef91202a91c9f3459007ac"),
			},
			scriptPubKey: hex2bytes("5120e4d810fd50586274face62b8a807eb9719cef49c04177cc6b76a9a4251d5450e"),
		},

		{
			internalPublicKey: hex2bytes("ee4fe085983462a184015d1f782d6a5f8b9c2b60130aff050ce221ecf3786592"),
			scriptTree: MastBranch{
				&MastLeaf{
					Version: constants.TaprootLeafVersionTapscript,
					Script:  hex2bytes("20387671353e273264c495656e27e39ba899ea8fee3bb69fb2a680e22093447d48ac"),
				},
				&MastLeaf{
					Version: 250,
					Script:  hex2bytes("06424950333431"),
				},
			},
			scriptPubKey: hex2bytes("5120712447206d7a5238acc7ff53fbe94a3b64539ad291c7cdbc490b7577e4b17df5"),
		},

		{
			internalPublicKey: hex2bytes("f9f400803e683727b14f463836e1e78e1c64417638aa066919291a225f0e8dd8"),
			scriptTree: MastBranch{
				&MastLeaf{
					Version: constants.TaprootLeafVersionTapscript,
					Script:  hex2bytes("2044b178d64c32c4a05cc4f4d1407268f764c940d20ce97abfd44db5c3592b72fdac"),
				},
				&MastLeaf{
					Version: constants.TaprootLeafVersionTapscript,
					Script:  hex2bytes("07546170726f6f74"),
				},
			},
			scriptPubKey: hex2bytes("512077e30a5522dd9f894c3f8b8bd4c4b2cf82ca7da8a3ea6a239655c39c050ab220"),
		},

		{ // Duplicate of the previous vector, but using the leaf hashes directly without the script preimages.
			internalPublicKey: hex2bytes("f9f400803e683727b14f463836e1e78e1c64417638aa066919291a225f0e8dd8"),
			scriptTree: MastBranch{
				MastLeafHash{
					0x64, 0x51, 0x2f, 0xec, 0xdb, 0x5a, 0xfa, 0x04, 0xf9, 0x88, 0x39, 0xb5, 0x0e, 0x6f, 0x0c, 0xb7,
					0xb1, 0xe5, 0x39, 0xbf, 0x6f, 0x20, 0x5f, 0x67, 0x93, 0x40, 0x83, 0xcd, 0xcc, 0x3c, 0x8d, 0x89,
				},
				MastLeafHash{
					0x2c, 0xb2, 0xb9, 0x0d, 0xaa, 0x54, 0x3b, 0x54, 0x41, 0x61, 0x53, 0x0c, 0x92, 0x5f, 0x28, 0x5b,
					0x06, 0x19, 0x69, 0x40, 0xd6, 0x08, 0x5c, 0xa9, 0x47, 0x4d, 0x41, 0xdc, 0x38, 0x22, 0xc5, 0xcb,
				},
			},
			scriptPubKey: hex2bytes("512077e30a5522dd9f894c3f8b8bd4c4b2cf82ca7da8a3ea6a239655c39c050ab220"),
		},

		{
			internalPublicKey: hex2bytes("e0dfe2300b0dd746a3f8674dfd4525623639042569d829c7f0eed9602d263e6f"),
			scriptTree: MastBranch{
				&MastLeaf{
					Version: constants.TaprootLeafVersionTapscript,
					Script:  hex2bytes("2072ea6adcf1d371dea8fba1035a09f3d24ed5a059799bae114084130ee5898e69ac"),
				},
				MastBranch{
					&MastLeaf{
						Version: constants.TaprootLeafVersionTapscript,
						Script:  hex2bytes("202352d137f2f3ab38d1eaa976758873377fa5ebb817372c71e2c542313d4abda8ac"),
					},
					&MastLeaf{
						Version: constants.TaprootLeafVersionTapscript,
						Script:  hex2bytes("207337c0dd4253cb86f2c43a2351aadd82cccb12a172cd120452b9bb8324f2186aac"),
					},
				},
			},
			scriptPubKey: hex2bytes("512091b64d5324723a985170e4dc5a0f84c041804f2cd12660fa5dec09fc21783605"),
		},

		{
			internalPublicKey: hex2bytes("55adf4e8967fbd2e29f20ac896e60c3b0f1d5b0efa9d34941b5958c7b0a0312d"),
			scriptTree: MastBranch{
				&MastLeaf{
					Version: constants.TaprootLeafVersionTapscript,
					Script:  hex2bytes("2071981521ad9fc9036687364118fb6ccd2035b96a423c59c5430e98310a11abe2ac"),
				},
				MastBranch{
					&MastLeaf{
						Version: constants.TaprootLeafVersionTapscript,
						Script:  hex2bytes("20d5094d2dbe9b76e2c245a2b89b6006888952e2faa6a149ae318d69e520617748ac"),
					},
					&MastLeaf{
						Version: constants.TaprootLeafVersionTapscript,
						Script:  hex2bytes("20c440b462ad48c7a77f94cd4532d8f2119dcebbd7c9764557e62726419b08ad4cac"),
					},
				},
			},
			scriptPubKey: hex2bytes("512075169f4001aa68f15bbed28b218df1d0a62cbbcf1188c6665110c293c907b831"),
		},
	}

	for _, fixture := range fixtures {
		scriptPubKey, err := MakeP2TR(fixture.internalPublicKey, fixture.scriptTree)
		if err != nil {
			t.Errorf("Failed to make P2TR address for key %x: %s", fixture.internalPublicKey, err)
			return
		}

		if !bytes.Equal(scriptPubKey, fixture.scriptPubKey) {
			t.Errorf("incorrect P2TR script\nWanted %x\nGot    %x", fixture.scriptPubKey, scriptPubKey)
			return
		}
	}
}

// Uses example from https://github.com/bitcoin-core/btcdeb/blob/1e7ff0f8da533bb8cf8ff6ebb4f160852cbb3685/doc/tapscript-example-with-tap.md
func TestMakeP2TR_btcdeb(t *testing.T) {
	// script1:
	//  144
	//  OP_CHECKSEQUENCEVERIFY
	//  OP_DROP
	//  <9997a497d964fc1a62885b05a51166a65a90df00492c8d7cf61d6accf54803be>
	//  OP_CHECKSIG
	alicePub := hex2bytes("9997a497d964fc1a62885b05a51166a65a90df00492c8d7cf61d6accf54803be")
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
	bobPub := hex2bytes("4edfcf9dfe6c0b5c83d1ab3f78d1b39a46ebac6798e08e19761f5ed89ec83c10")
	hashed := hex2bytes("6c60f404f8167a38fc70eaf8aa17ac351023bef86bcb9d1086a19afe95bd5333")

	script2 := append([]byte{constants.OP_SHA256}, PushData(hashed)...)
	script2 = append(script2, constants.OP_EQUALVERIFY)
	script2 = append(script2, PushData(bobPub)...)
	script2 = append(script2, constants.OP_CHECKSIG)

	mast := MastBranch{
		&MastLeaf{
			Version: constants.TaprootLeafVersionTapscript,
			Script:  script1,
		},
		&MastLeaf{
			Version: constants.TaprootLeafVersionTapscript,
			Script:  script2,
		},
	}

	mastRootHash := mast.Hash()
	expectedRootHash := hex2bytes("41646f8c1fe2a96ddad7f5471bc4fee7da98794ef8c45a4f4fc6a559d60c9f6b")
	if !bytes.Equal(mastRootHash[:], expectedRootHash) {
		t.Errorf("MAST root hash does not match\nWanted %x\nGot    %x", expectedRootHash, mastRootHash)
		return
	}

	internalPublicKey := hex2bytes("f30544d6009c8d8d94f5d030b2e844b1a3ca036255161c479db1cca5b374dd1c")
	scriptPubKey, err := MakeP2TR(internalPublicKey, mast)
	if err != nil {
		t.Errorf("failed to build P2TR output script: %s", err)
		return
	}

	expectedScript := hex2bytes("5120a5ba0871796eb49fb4caa6bf78e675b9455e2d66e751676420f8381d5dda8951")
	if !bytes.Equal(scriptPubKey, expectedScript) {
		t.Errorf("P2TR script does not match\nWanted %x\nGot    %x", expectedScript, scriptPubKey)
		return
	}
}
