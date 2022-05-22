package blocks

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

func TestBlock(t *testing.T) {
	type Fixture struct {
		BlockNumber   int      `json:"blockNumber"`
		RawHex        string   `json:"raw"`
		HeaderHashHex string   `json:"hash"`
		Txids         []string `json:"txids"`
		Weight        int      `json:"weight"`
	}

	var fixtures []*Fixture

	file, err := os.Open("fixtures.json")
	if err != nil {
		t.Errorf("failed to open block fixtures file: %s", err)
		return
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(&fixtures); err != nil {
		t.Errorf("failed to decode block fixtures JSON: %s", err)
		return
	}

	for _, fixture := range fixtures {
		blockBytes, err := hex.DecodeString(fixture.RawHex)
		if err != nil {
			t.Errorf("failed to decode block fixture hex: %s", err)
			return
		}

		block, err := FromReader(bytes.NewReader(blockBytes))
		if err != nil {
			t.Errorf("failed to decode block %d: %s", fixture.BlockNumber, err)
			return
		}

		headerHash, err := block.Header.Hash()
		if err != nil {
			t.Errorf("failed to hash block %d header: %s", fixture.BlockNumber, err)
			return
		}

		headerHashHex := hex.EncodeToString(headerHash[:])
		if headerHashHex != fixture.HeaderHashHex {
			t.Errorf(
				"block %d header hash does not match\nWanted %s\nGot    %s",
				fixture.BlockNumber,
				fixture.HeaderHashHex,
				headerHashHex,
			)
			return
		}

		if len(block.Transactions) != len(fixture.Txids) {
			t.Errorf(
				"block %d has wrong number of transactions\nWanted %d\nGot    %d",
				fixture.BlockNumber,
				len(fixture.Txids),
				len(block.Transactions),
			)
			return
		}

		for i := 0; i < len(block.Transactions); i++ {
			txid, err := block.Transactions[i].Id(false)
			if err != nil {
				t.Errorf("in block %d, failed to hash transaction %d: %s", fixture.BlockNumber, i, err)
				return
			}

			expected := fixture.Txids[i]
			if txid != expected {
				t.Errorf(
					"in block %d, TXID %d does not match\nWanted %s\nGot    %s",
					fixture.BlockNumber,
					i,
					expected,
					txid,
				)
				return
			}
		}

		serializedBlock := block.Bytes()
		if !bytes.Equal(serializedBlock, blockBytes) {
			t.Errorf("block %d did not re-serialize correctly", fixture.BlockNumber)
			return
		}

		if blockSize := block.Size(); blockSize != len(serializedBlock) {
			t.Errorf("block Size() did not match expected\nWanted %d\nGot    %d", len(serializedBlock), blockSize)
			return
		}

		if blockWeight := block.WeightUnits(); blockWeight != fixture.Weight {
			t.Errorf("block WeightUnits() did not match expected\nWanted %d\nGot    %d", fixture.Weight, blockWeight)
			return
		}
	}
}
