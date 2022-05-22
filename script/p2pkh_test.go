package script

import (
	"bytes"
	"testing"

	"github.com/kklash/bitcoinlib/bhash"
)

func TestMakeDecodeP2PKH(t *testing.T) {
	type Fixture struct {
		publicKey []byte
		script    []byte
	}

	fixtures := []Fixture{
		Fixture{
			hex2bytes("0302d4d85ec9610122e1bae04a770755ae15fd8b5c8c305b15ce09e4eae52b90fd"),
			hex2bytes("76a91491c79c05a31adead59033ebf47acab299b4cdba488ac"),
		},
		Fixture{
			hex2bytes("030cbb649a9b35c96ec6649624cb2f4c1bc44db90c9d2995c7219ab527714dfdc8"),
			hex2bytes("76a9145c297fca6dd69f40fd6435b4c4b626739ab8aed788ac"),
		},
		Fixture{
			hex2bytes("028e4326aea645a4e5caad53f085e273596686bb2dd755629e9abc5161ff356cbd"),
			hex2bytes("76a914a4e2bb2b54534f57f87d607b69f681c96111db9588ac"),
		},
	}

	for _, fixture := range fixtures {
		hash := bhash.Hash160(fixture.publicKey)
		scriptPubKey := MakeP2PKHFromHash(hash)
		if !bytes.Equal(scriptPubKey, fixture.script) {
			t.Errorf("script pub key does not match:\n wanted %x\n got %x", fixture.script, scriptPubKey)
			continue
		}

		scriptPubKeyFromPublicKey, err := MakeP2PKHFromPublicKey(fixture.publicKey)
		if err != nil {
			t.Errorf("failed to make script pub key from public key")
			continue
		}

		if !bytes.Equal(scriptPubKeyFromPublicKey, scriptPubKey) {
			t.Errorf("script pub key from public key does not match:\n wanted %x\n got %x", scriptPubKey, scriptPubKeyFromPublicKey)
			continue
		}

		decodedHash, err := DecodeP2PKH(scriptPubKey)
		if err != nil {
			t.Errorf("failed to decode script pub key: %s", err)
			continue
		}

		if decodedHash != hash {
			t.Errorf("decoded pk-hash does not match fixture\n wanted %x\n got %x", hash, decodedHash)
		}
	}
}

func TestIsP2PKH(t *testing.T) {
	valid := [][]byte{
		hex2bytes("76a914c41c836560406c6169537f9dc9520184879f03e288ac"),
		hex2bytes("76a914e612574a705f1992f9cf86b4575a779b39ca074088ac"),
		hex2bytes("76a914fa18be26c850da95b7827c2b1474d20c7849ab4288ac"),
	}

	invalid := [][]byte{
		nil,
		hex2bytes("deadbeef"),
		hex2bytes("76a914fa18be26c850da95b7827c2b1474d20c7849ab4288ac01"),
		hex2bytes("0076a914943f55fc74de5ea179f5f05ecb7396a1bf50c44088ac"),
	}

	for _, scriptPubKey := range valid {
		if !IsP2PKH(scriptPubKey) {
			t.Errorf("failed to recognize P2PKH script: %x", scriptPubKey)
		}
	}

	for _, scriptPubKey := range invalid {
		if IsP2PKH(scriptPubKey) {
			t.Errorf("detected invalid script as P2PKH: %x", scriptPubKey)
		}
	}
}

func TestRedeemP2PKH(t *testing.T) {
	sig := hex2bytes("3045022100854e2b57a878df26971911135ab9bb851e4de4560d29fcd57a3403d92a27485a02207378f098143caa6ce4d31fcba3f753f33f4a3ecc03c51ba465e8dd6f44e053f301")
	pk := hex2bytes("02290559e16eb8f6ba4afbe613cfe5f1e1c000c345ba3fe192938f4be4bba385ae")

	expected := hex2bytes("483045022100854e2b57a878df26971911135ab9bb851e4de4560d29fcd57a3403d92a27485a02207378f098143caa6ce4d31fcba3f753f33f4a3ecc03c51ba465e8dd6f44e053f3012102290559e16eb8f6ba4afbe613cfe5f1e1c000c345ba3fe192938f4be4bba385ae")

	actual := RedeemP2PKH(sig, pk)
	if !bytes.Equal(actual, expected) {
		t.Errorf("P2PKH redeem script did not build correctly\nWanted %x\nGot    %x", expected, actual)
	}
}
