package blockscan

import (
	"context"
	"encoding/hex"
	"testing"
)

func TestBlockScanner_StreamBlocksUnordered(t *testing.T) {
	t.Skip()

	scanner, err := createScanner()
	if err != nil {
		t.Errorf("failed to create scanner: %s", err)
		return
	}

	var fromHeight uint32 = 100
	var toHeight uint32 = 500

	nextBlock := scanner.StreamBlocksUnordered(context.TODO(), fromHeight, toHeight, 20)

	var blocksCounted int = 0
	for ; ; blocksCounted++ {
		actualBlock, err := nextBlock()
		if err != nil {
			t.Errorf("Failed to get unordered block: %s", err)
			return
		}

		if actualBlock == nil {
			break
		}

		hash, err := actualBlock.Header.Hash()
		if err != nil {
			t.Errorf("failed to hash block: %s", err)
			return
		}

		var expectedBlock struct {
			Height uint32
		}
		err = scanner.Connection.RequestSetResult(
			&expectedBlock,
			"getblock",
			hex.EncodeToString(hash[:]),
		)
		if err != nil {
			t.Errorf("failed to get expected block %x: %s", hash[:], err)
			return
		}

		if expectedBlock.Height < fromHeight || expectedBlock.Height > toHeight {
			t.Errorf("fetched block of unexpected height: %d", expectedBlock.Height)
			return
		}
	}

	expectedBlocksCounted := int(toHeight-fromHeight) + 1
	if blocksCounted != expectedBlocksCounted {
		t.Errorf("scanned wrong number of blocks\nwanted %d\ngot    %d", expectedBlocksCounted, blocksCounted)
		return
	}
}
