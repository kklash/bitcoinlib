package bip39

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/kklash/bitcoinlib/bip32"
	"github.com/kklash/bitcoinlib/constants"
)

const TestVectorPassphrase = "TREZOR"

func TestBip39(t *testing.T) {
	test := func(entHex, expectedMnemonic, expectedSeedHex, expectedXpriv string) {
		entropy, _ := hex.DecodeString(entHex)
		mnemonic, err := Encode(entropy)
		if err != nil {
			t.Errorf("Failed to convert data to mnemonic: %s\nError: %s", expectedMnemonic, err)
			return
		}

		if mnemonic != expectedMnemonic {
			t.Errorf("Mnemonic does not match expected\nWanted %s\nGot    %s", expectedMnemonic, mnemonic)
			return
		}

		decoded, err := Decode(mnemonic)
		if err != nil {
			t.Errorf("Failed to decode mnemonic: %s\nError: %s", mnemonic, err)
			return
		}

		if decodedHex := hex.EncodeToString(decoded); decodedHex != entHex {
			t.Errorf("Decoded entropy does not match\nWanted %s\nGot    %s", entHex, decodedHex)
			return
		}

		seed := DeriveSeed(mnemonic, TestVectorPassphrase)
		if hex.EncodeToString(seed) != expectedSeedHex {
			t.Errorf("Derived seed does not match fixture\nWanted %s\nGot    %x", expectedSeedHex, seed)
			return
		}

		masterKey, chainCode, err := bip32.GenerateMasterKey(seed)
		if err != nil {
			t.Errorf("Failed to derive master key: %s\nError: %s", expectedXpriv, err)
			return
		}

		xpriv := bip32.SerializePrivate(masterKey, chainCode, nil, 0, 0, constants.BitcoinNetwork.ExtendedPrivate)
		if xpriv != expectedXpriv {
			t.Errorf("BIP32 derived master key does not match fixture\nWanted %s\nGot    %s", expectedXpriv, xpriv)
			return
		}
	}

	var fixtures [][]string
	data, err := ioutil.ReadFile("test_vectors.json")
	if err != nil {
		t.Errorf("Failed to read test vectors: %s", err)
		return
	}

	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Errorf("Failed to decode test vectors: %s", err)
		return
	}

	for _, fixture := range fixtures {
		test(fixture[0], fixture[1], fixture[2], fixture[3])
	}
}
