package blockscan

import (
	"context"
	"encoding/base64"
	"testing"
)

func TestBlockScanner_StreamBlocks(t *testing.T) {
	t.Skip()

	scanner, err := createScanner()
	if err != nil {
		t.Errorf("failed to create scanner: %s", err)
		return
	}

	var fromHeight uint32 = 50000
	var toHeight uint32 = 50600

	nextBlock := scanner.StreamBlocks(context.TODO(), fromHeight, toHeight, 20)

	var height uint32 = fromHeight
	for ; ; height++ {
		block, err := nextBlock()
		if err != nil {
			t.Errorf("failed to get block: %s", err)
			return
		}

		if block == nil {
			height -= 1
			break
		}

		expectedBlock, err := scanner.GetBlockByHeight(height)
		if err != nil {
			t.Errorf("failed to get expected block: %s", err)
			return
		}

		blockBase64 := base64.StdEncoding.EncodeToString(block.Bytes())
		expectedBlockBase64 := base64.StdEncoding.EncodeToString(expectedBlock.Bytes())

		if blockBase64 != expectedBlockBase64 {
			t.Errorf("block at height %d does not match\nWanted %s\nGot    %s", height, expectedBlockBase64, blockBase64)
			return
		}
	}

	if height != toHeight {
		t.Errorf("did not scan expected number of blocks")
		return
	}
}
