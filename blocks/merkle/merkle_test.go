package merkle

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

type Fixture struct {
	Comment    string   `json:"comment"`
	MerkleRoot string   `json:"merkleroot"`
	TxidsHex   []string `json:"tx"`
	Txids      [][32]byte
}

func getFixtures() ([]*Fixture, error) {
	fixturesBytes, err := ioutil.ReadFile("fixtures.json")
	if err != nil {
		return nil, fmt.Errorf("Failed to read from fixtures file: %s", err)
	}

	var fixtures []*Fixture
	if err := json.Unmarshal(fixturesBytes, &fixtures); err != nil {
		return nil, fmt.Errorf("Failed to decode block fixtures: %s", err)
	}

	for _, fixture := range fixtures {
		fixture.Txids = make([][32]byte, len(fixture.TxidsHex))
		for i, txidHex := range fixture.TxidsHex {
			txid, err := hex.DecodeString(txidHex)
			if err != nil {
				return nil, fmt.Errorf("failed to decode TXID: %s - %s", txidHex, err)
			}

			copy(fixture.Txids[i][:], txid)
		}
	}

	return fixtures, nil
}

func TestMerkleTree(t *testing.T) {
	fixtures, err := getFixtures()
	if err != nil {
		t.Error(err)
		return
	}

	for _, fixture := range fixtures {
		merkle := MerkleRootHash(fixture.Txids)
		merkleHex := hex.EncodeToString(merkle[:])

		if merkleHex != fixture.MerkleRoot {
			t.Errorf("merkle root hash does not match fixture\nWanted %s\nGot    %s", fixture.MerkleRoot, merkleHex)
			return
		}
	}
}
