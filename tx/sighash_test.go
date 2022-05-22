package tx

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/kklash/bitcoinlib/common"
)

func mustHex(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return decoded
}

func TestSigHash(t *testing.T) {
	sigHashTestsJSON, err := ioutil.ReadFile("sighash.json")
	if err != nil {
		t.Errorf("ERROR: %s", err)
		return
	}

	var sigHashTests [][]interface{}

	if err := json.Unmarshal(sigHashTestsJSON, &sigHashTests); err != nil {
		t.Errorf("ERROR: %s", err)
		return
	}

	wins := 0

	for _, fixture := range sigHashTests[1:] {
		rawTx := mustHex(fixture[0].(string))
		prevOutScript := mustHex(fixture[1].(string))
		nInput := int(fixture[2].(float64))
		hashType := uint32(fixture[3].(float64))
		var expectedSigHash [32]byte
		copy(expectedSigHash[:], common.ReverseBytes(mustHex(fixture[4].(string))))

		tx, err := FromBytes(rawTx)
		if err != nil {
			t.Errorf("ERROR: %s", err)
			return
		}

		sigHash, err := tx.SignatureHashForInput(nInput, prevOutScript, hashType)
		if err != nil {
			t.Errorf("ERROR: %s", err)
			return
		}

		if sigHash == expectedSigHash {
			wins += 1
		} else {
			t.Errorf("FAIL: mismatch - %x != %x\n", sigHash, expectedSigHash)

			if bytes.Contains(prevOutScript, []byte{0xab}) {
				t.Errorf("  Probably due to code-separator: %x\n", prevOutScript)
			}
		}
	}
}

// These fixtures taken directly from BIP-0143
func TestSigHashWitness(t *testing.T) {
	type Fixture struct {
		TxHex string `json:"tx"`
		Input struct {
			Index           int    `json:"index"`
			Value           uint64 `json:"value"`
			ScriptPubKeyHex string `json:"scriptPubKey"`
		} `json:"input"`
		SigHashType uint32 `json:"sighashType"`
		SigHashHex  string `json:"sighash"`
	}

	sigHashTestsJSON, err := ioutil.ReadFile("sighash_witness.json")
	if err != nil {
		t.Errorf("ERROR: %s", err)
		return
	}

	var sigHashTests []*Fixture
	if err := json.Unmarshal(sigHashTestsJSON, &sigHashTests); err != nil {
		t.Errorf("ERROR: %s", err)
		return
	}

	for _, fixture := range sigHashTests {
		txn, err := FromBytes(mustHex(fixture.TxHex))
		if err != nil {
			t.Errorf("ERROR: %s", err)
			return
		}

		sigHash, err := txn.SignatureHashForWitnessInput(
			fixture.Input.Index,
			mustHex(fixture.Input.ScriptPubKeyHex),
			fixture.SigHashType,
			fixture.Input.Value,
		)
		if err != nil {
			t.Errorf("ERROR: %s", err)
			return
		}

		sigHashHex := fmt.Sprintf("%x", sigHash)
		if sigHashHex != fixture.SigHashHex {
			t.Errorf("Signature hash does not match\nWanted %s\nGot    %s", fixture.SigHashHex, sigHashHex)
			continue
		}
	}
}
