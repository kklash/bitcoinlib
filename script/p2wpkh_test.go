package script

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/kklash/bitcoinlib/bhash"
)

func hex2bytes(h string) []byte {
	buf, _ := hex.DecodeString(h)
	return buf
}

func TestMakeP2WPKH(t *testing.T) {
	type Fixture struct {
		publicKey []byte
		script    []byte
	}

	fixtures := []Fixture{
		Fixture{
			hex2bytes("026d49008eb40f28c497e32a4c2220d893a9b25f7594f7081193c3a3ee79b77f86"),
			hex2bytes("0014d8261b1642dc77b064caf9086cc3ed18ea7cd1a4"),
		},
		Fixture{
			hex2bytes("03439c7d975fb0580df6760225aeb9572e4469ae2c8191ff7438050ec2b1a00c78"),
			hex2bytes("0014b21658ea3960030b4dc4634f51ecae0505daf2e6"),
		},
		Fixture{
			hex2bytes("02dd9284f10e2ecfc4852364c6a38f49ba2f1a76395e4b038bba3f7297bed986c8"),
			hex2bytes("0014ce6a28589e056b0bdd67c464033677a7ac35ce05"),
		},
	}

	for _, fixture := range fixtures {
		hash := bhash.Hash160(fixture.publicKey)
		scriptPubKey := MakeP2WPKHFromHash(hash)
		if !bytes.Equal(scriptPubKey, fixture.script) {
			t.Errorf("script pub key does not match:\n wanted %x\n got %x", fixture.script, scriptPubKey)
			continue
		}

		scriptPubKeyFromPublicKey, err := MakeP2WPKHFromPublicKey(fixture.publicKey)
		if err != nil {
			t.Errorf("failed to make script pub key from public key")
			continue
		}

		if !bytes.Equal(scriptPubKeyFromPublicKey, scriptPubKey) {
			t.Errorf("script pub key from public key does not match:\n wanted %x\n got %x", scriptPubKey, scriptPubKeyFromPublicKey)
			continue
		}

		decodedHash, err := DecodeP2WPKH(scriptPubKey)
		if err != nil {
			t.Errorf("failed to decode script pub key: %s", err)
			continue
		}

		if decodedHash != hash {
			t.Errorf("decoded pk-hash does not match fixture\n wanted %x\n got %x", hash, decodedHash)
		}
	}
}

func TestIsP2WPKH(t *testing.T) {
	valid := [][]byte{
		hex2bytes("0014098a0d10c22aa5bf136c5d78db397e2e27772189"),
		hex2bytes("0014c67e4ce743057b38b479649074f0497134231e66"),
		hex2bytes("0014d0bf1e07009a235ce49a460f7d1436da94fbe6e0"),
	}

	invalid := [][]byte{
		hex2bytes("deadbeef"),
		nil,
		hex2bytes("0012b76272198beb0c95cd26904ffa589250b3ff0494"),
	}

	for _, scriptPubKey := range valid {
		if !IsP2WPKH(scriptPubKey) {
			t.Errorf("failed to recognize P2WPKH script: %x", scriptPubKey)
		}
	}

	for _, scriptPubKey := range invalid {
		if IsP2WPKH(scriptPubKey) {
			t.Errorf("detected invalid script as P2WPKH: %x", scriptPubKey)
		}
	}
}
