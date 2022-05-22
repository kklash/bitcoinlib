package script

import (
	"bytes"
	"testing"

	"github.com/kklash/bitcoinlib/bhash"
)

func TestMakeP2SH(t *testing.T) {
	type Fixture struct {
		inputScript  []byte
		scriptPubKey []byte
	}

	fixtures := []Fixture{
		Fixture{
			hex2bytes("00147b5c838604f821aa542e53b53610cb821c6b27c1"), // P2WPKH
			hex2bytes("a9146505c5df8a4dcde5d6ef68cb4c15efa19c01991c87"),
		},
		Fixture{
			hex2bytes("76a91443f2b69a8c542c6c2f5646a74ecd39df0961cf3288ac"), // P2PKH
			hex2bytes("a91490211802f3e53549e518ec22f812f13c4954930887"),
		},
		Fixture{
			hex2bytes("0020acf201c4c0fadcf9dc21696a0046034732b743eac97a2826beb9877ef9829fc8"), // P2WSH
			hex2bytes("a9149fa3a434877e0f36d2e78223d448ea549caac06c87"),
		},
	}

	for _, fixture := range fixtures {
		hash := bhash.Hash160(fixture.inputScript)
		scriptPubKey := MakeP2SHFromHash(hash)

		if !bytes.Equal(scriptPubKey, fixture.scriptPubKey) {
			t.Errorf("script pub key does not match:\n wanted %x\n got    %x", fixture.scriptPubKey, scriptPubKey)
			continue
		}

		scriptPubKeyFromScript := MakeP2SHFromScript(fixture.inputScript)
		if !bytes.Equal(scriptPubKey, scriptPubKeyFromScript) {
			t.Errorf("script pub key does not match:\n wanted %x\n got    %x", scriptPubKey, scriptPubKeyFromScript)
			continue
		}

		decodedHash, err := DecodeP2SH(scriptPubKey)
		if err != nil {
			t.Errorf("failed to decode script pub key: %s", err)
			continue
		}

		if decodedHash != hash {
			t.Errorf("decoded script-hash does not match fixture\n wanted %x\n got    %x", hash, decodedHash)
		}
	}
}

func TestIsP2SH(t *testing.T) {
	valid := [][]byte{
		hex2bytes("a914aefdfa5219b0ecda7f0f75a2221111769dc848c187"),
		hex2bytes("a9142882270c08b9ebd20117526a0f53f3b2ecd01de987"),
		hex2bytes("a91424f57aeedf0e90d198ca1f5932d2de4ebf659b3387"),
	}

	invalid := [][]byte{
		nil,
		hex2bytes("feeddeeb"),
		hex2bytes("a900140f0e3d67a375c09b2c91811807f417059681c10c87"),
		hex2bytes("a9141bfd81193c2f21530fe61fbe406aebfb4b5510288700"),
	}

	for _, scriptPubKey := range valid {
		if !IsP2SH(scriptPubKey) {
			t.Errorf("failed to recognize P2SH script: %x", scriptPubKey)
		}
	}

	for _, scriptPubKey := range invalid {
		if IsP2SH(scriptPubKey) {
			t.Errorf("detected invalid script as P2SH: %x", scriptPubKey)
		}
	}
}

func TestRedeemP2SH(t *testing.T) {
	pk1 := hex2bytes("022afc20bf379bc96a2f4e9e63ffceb8652b2b6a097f63fbee6ecec2a49a48010e")
	pk2 := hex2bytes("03a767c7221e9f15f870f1ad9311f5ab937d79fcaeee15bb2c722bca515581b4c0")
	multiSigScriptPubKey := MakeP2MS(1, pk1, pk2)
	underlyingRedemptionScript := RedeemP2MS(
		hex2bytes("3046022100a07b2821f96658c938fa9c68950af0e69f3b2ce5f8258b3a6ad254d4bc73e11e022100e82fab8df3f7e7a28e91b3609f91e8ebf663af3a4dc2fd2abd954301a5da67e701"),
	)

	expectedScriptSig := hex2bytes("00493046022100a07b2821f96658c938fa9c68950af0e69f3b2ce5f8258b3a6ad254d4bc73e11e022100e82fab8df3f7e7a28e91b3609f91e8ebf663af3a4dc2fd2abd954301a5da67e701475121022afc20bf379bc96a2f4e9e63ffceb8652b2b6a097f63fbee6ecec2a49a48010e2103a767c7221e9f15f870f1ad9311f5ab937d79fcaeee15bb2c722bca515581b4c052ae")

	actualScriptSig := RedeemP2SH(multiSigScriptPubKey, underlyingRedemptionScript)

	if !bytes.Equal(actualScriptSig, expectedScriptSig) {
		t.Errorf("P2SH redeem script did not build correctly\nWanted %x\nGot    %x", expectedScriptSig, actualScriptSig)
	}
}
