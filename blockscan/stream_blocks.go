package blockscan

import (
	"context"
	"fmt"

	"github.com/kklash/bitcoinlib/blocks"
)

func (scanner *BlockScanner) streamBlocks(
	ctx context.Context,
	fromHeight, toHeight uint32,
	parallelism uint32,
) (chan *blocks.Block, chan error) {
	cancelCtx, cancel := context.WithCancel(ctx)

	orderedBlockQueue := make(chan *blocks.Block)
	wrappedErrorQueue := make(chan error)

	go func() {
		defer func() {
			cancel()
			close(orderedBlockQueue)
			close(wrappedErrorQueue)
		}()

		firstBlock, err := scanner.GetBlockByHeight(fromHeight)
		if err != nil {
			wrappedErrorQueue <- err
			return
		}

		latestBlockHash, err := firstBlock.Header.Hash()
		if err != nil {
			wrappedErrorQueue <- err
			return
		}

		blockQueue, errorQueue := scanner.streamBlocksUnordered(
			cancelCtx,
			fromHeight+1,
			toHeight,
			parallelism,
		)

		// creates a sort of forward-looking linked list of blocks,
		// mapping the hash of each block to the next block in the chain.
		blocksByPrevHash := make(map[[32]byte]*blocks.Block)

		orderedBlockQueue <- firstBlock

		var channelsClosed bool

		for currentHeight := fromHeight; currentHeight < toHeight; {
			select {
			case <-ctx.Done():
				return

			case block, more := <-blockQueue:
				if !more {
					channelsClosed = true
					break
				}

				blocksByPrevHash[block.Header.PreviousHeaderHash] = block

			case err, more := <-errorQueue:
				if !more {
					channelsClosed = true
					break
				}

				wrappedErrorQueue <- err
				return
			}

			for {
				nextBlock, ok := blocksByPrevHash[latestBlockHash]
				if ok {
					delete(blocksByPrevHash, latestBlockHash)
					orderedBlockQueue <- nextBlock
					currentHeight += 1

					latestBlockHash, err = nextBlock.Header.Hash()
					if err != nil {
						wrappedErrorQueue <- err
						return
					}
				} else {
					break
				}
			}

			if channelsClosed {
				wrappedErrorQueue <- fmt.Errorf(
					"blockchain link broken at block hash %x - no next block found",
					latestBlockHash,
				)
				return
			}
		}
	}()

	return orderedBlockQueue, wrappedErrorQueue
}

// StreamBlocks streams blocks in ascending height order, beginning at startBlockHeight, and
// ending before endBlockHeight (exclusive range). It returns a NextBlockFunc which will return the
// next block with each sequential call. An in-progress stream can be canceled by the BlockScanner's
// context.
func (scanner *BlockScanner) StreamBlocks(
	ctx context.Context,
	startBlockHeight, endBlockHeight uint32,
	parallelism uint32,
) NextBlockFunc {
	blockQueue, errorQueue := scanner.streamBlocks(
		ctx,
		startBlockHeight,
		endBlockHeight,
		parallelism,
	)

	return newNextBlockFunc(blockQueue, errorQueue)
}
