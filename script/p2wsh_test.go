package script

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/kklash/bitcoinlib/bhash"
)

func TestMakeP2WSH(t *testing.T) {
	type Fixture struct {
		inputScript  []byte
		scriptPubKey []byte
	}

	fixtures := []Fixture{
		Fixture{
			hex2bytes("00147b5c838604f821aa542e53b53610cb821c6b27c1"), // P2WPKH
			hex2bytes("00203298ddb78025128d6424a089f4aa233917f9c8fcfcf7949840222e0bb9ac1adc"),
		},
		Fixture{
			hex2bytes("0400e10b5eb175a914ea5962e5e1977db9eeb08bd521267dde5242392388ac"), // locktime script
			hex2bytes("0020019b20e82a79307b33ac3692e30e4065dbbc6c5643d1333a77e13a3d457d1730"),
		},
		Fixture{
			hex2bytes("5221037070dc9c046ce2cd18b5847899a016fafdd94f61c4032a91c8b33111c43f3bce21022bf784c5e4fee1dc7944fd723356304b4b307689fb9aa1a56c95177d6811d05321025257a485751eb4e9a264d9bc8632b9e5c25c24f6beffc159ce7b08e110ebe4d653ae"), // multisig
			hex2bytes("00202696ac1fed6f5e756bb6a36ac912d996fdf602c9705442c8a61bf240fc583872"),
		},
	}

	for _, fixture := range fixtures {
		hash := bhash.Sha256(fixture.inputScript)
		scriptPubKey := MakeP2WSHFromHash(hash)

		if !bytes.Equal(scriptPubKey, fixture.scriptPubKey) {
			t.Errorf("script pub key does not match:\n wanted %x\n got %x", fixture.scriptPubKey, scriptPubKey)
			continue
		}

		scriptPubKeyFromScript := MakeP2WSHFromScript(fixture.inputScript)
		if !bytes.Equal(scriptPubKey, scriptPubKeyFromScript) {
			t.Errorf("script pub key does not match:\n wanted %x\n got    %x", scriptPubKey, scriptPubKeyFromScript)
			continue
		}

		decodedHash, err := DecodeP2WSH(scriptPubKey)
		if err != nil {
			t.Errorf("failed to decode script pub key: %s", err)
			continue
		}

		if decodedHash != hash {
			t.Errorf("decoded script-hash does not match fixture\n wanted %x\n got %x", hash, decodedHash)
		}
	}
}

func TestIsP2WSH(t *testing.T) {
	valid := [][]byte{
		hex2bytes("0020019b20e82a79307b33ac3692e30e4065dbbc6c5643d1333a77e13a3d457d1730"),
		hex2bytes("00205f78c33274e43fa9de5659265c1d917e25c03722dcb0b8d27db8d5feaa813953"),
		hex2bytes("0020a745dc51b20631cd78971b67274bc1b1cb06fb800a408f4184fb2bf00dfe9b9d"),
	}

	invalid := [][]byte{
		hex2bytes("deadbeef"),
		hex2bytes("00205221037070dc9c046ce2cd18b5847899a016fafdd94f61c4032a91c8b33111c43f3bce21022b"),
		hex2bytes("0021e60d79537e948624ee19fd91edb17035046a0d6888a78ecfcaa4eb11fe4afa9b"),
	}

	for _, scriptPubKey := range valid {
		if !IsP2WSH(scriptPubKey) {
			t.Errorf("failed to recognize P2WSH script: %x", scriptPubKey)
		}
	}

	for _, scriptPubKey := range invalid {
		if IsP2WSH(scriptPubKey) {
			t.Errorf("detected invalid script as P2WSH: %x", scriptPubKey)
		}
	}
}

func TestWitnessP2WSH(t *testing.T) {
	pk1 := hex2bytes("03b5999c30e64da515bb1990d2c10fcee296371ed98c14fc84bc258edd5ee7f1e2")
	pk2 := hex2bytes("02f28b7e07eb11a84c65db14939595f173a289987bdfb06150d5ff885bbd5827ba")
	pk3 := hex2bytes("03d0cfa728d0a3b71afbc27f1d96a7da019a2c978aa0c74b4a137180d10c2b366a")

	multisigScriptPubKey := MakeP2MS(2, pk1, pk2, pk3)

	sig1 := hex2bytes("3043021f2723ef0bc37f345b7c7388bd4008963662af2a8c9d09f93dc3b113efc9cf8a022009f9ab57bbe76d9b9d3fcde7c745d1c57d741adb7c36157c872887a9bb329df301")
	sig2 := hex2bytes("3045022100f9cafaa1a2d2a248a47eda156e21b79110b82b89215ccae9c347ad8f713da51202201e8c40509b33fa18c10307c71cc341faa89cd6d3d6cdf0c1d80cc26ea21e8d0b01")
	multisigScriptSig := append(RedeemP2MS(sig1, sig2))

	w, err := WitnessP2WSH(multisigScriptPubKey, multisigScriptSig)
	if err != nil {
		t.Errorf("Failed to create witness for P2WSH: %s", err)
		return
	}
	expectedWitness := [][]byte{
		{},
		sig1,
		sig2,
		multisigScriptPubKey,
	}

	if !reflect.DeepEqual(w, expectedWitness) {
		t.Errorf("failed to build expected P2WSH witness")
		return
	}
}
