package blockscan

import (
	"context"

	"github.com/kklash/bitcoinlib/unspent"
)

// UpdateUtxos scans the given block range for UTXOs belonging to the given set of scriptPubKeys,
// starting at startBlockHeight and ending before endBlockHeight (exclusive range). It updates
// a given set of unspent outputs which should be up to date as of startBlockHeight-1. If a non-nil
// progressCallback is passed, it is called after each block is scanned successfully. Block
// streaming is cancelled if the given Context ctx is cancelled or if an error is encountered.
func (scanner *BlockScanner) UpdateUtxos(
	ctx context.Context,
	utxos *unspent.OutputSet,
	scriptPubKeys [][]byte,
	startBlockHeight, endBlockHeight uint32,
	parallelism uint32,
	progressCallback func(height uint32),
) error {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	getNextBlock := scanner.StreamBlocks(cancelCtx, startBlockHeight, endBlockHeight, parallelism)

	for scanBlockHeight := startBlockHeight; scanBlockHeight <= endBlockHeight; scanBlockHeight++ {
		block, err := getNextBlock()
		if err != nil {
			return err
		}

		if err := utxos.UpdateFromBlock(block, scriptPubKeys); err != nil {
			return err
		}

		if progressCallback != nil {
			progressCallback(scanBlockHeight)
		}
	}

	return nil
}
