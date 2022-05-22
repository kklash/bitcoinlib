package blockscan

import (
	"context"
	"errors"
	"sync"

	"github.com/kklash/bitcoinlib/blocks"
)

// ErrInvalidBlockHeight is returned by BlockScanner's streaming methods if
// the starting block height is higher than the ending block height.
var ErrInvalidBlockHeight = errors.New("invalid block height specified")

func (scanner *BlockScanner) streamBlocksUnordered(
	ctx context.Context,
	fromHeight, toHeight uint32,
	parallelism uint32,
) (chan *blocks.Block, chan error) {
	if toHeight < fromHeight {
		panic(ErrInvalidBlockHeight)
	}

	blockQueue := make(chan *blocks.Block)
	errorQueue := make(chan error)

	var wg sync.WaitGroup

	wg.Add(int(parallelism))

	for i := uint32(0); i < parallelism; i++ {
		go func(i uint32) {
			defer wg.Done()

			blocksFetched := uint32(0)

			for {
				nBlockToGet := fromHeight + blocksFetched*parallelism + i
				if nBlockToGet > toHeight {
					return
				}

				block, err := scanner.GetBlockByHeight(nBlockToGet)
				if err != nil {
					errorQueue <- err
					return
				}

				select {
				case <-ctx.Done():
					return
				case blockQueue <- block:
					blocksFetched += 1
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(blockQueue)
		close(errorQueue)
	}()

	return blockQueue, errorQueue
}

// StreamBlocksUnordered streams blocks as quickly as possible without synchronizing to order them by height.
// Streaming begins at startBlockHeight, and ends before endBlockHeight (exclusive range). It returns a
// NextBlockFunc which will return a new block in the range [startBlockHeight, endBlockHeight] with each
// sequential call. An in-progress stream can be canceled by the BlockScanner's context.
func (scanner *BlockScanner) StreamBlocksUnordered(
	ctx context.Context,
	startBlockHeight, endBlockHeight uint32,
	parallelism uint32,
) NextBlockFunc {
	blockQueue, errorQueue := scanner.streamBlocksUnordered(
		ctx,
		startBlockHeight,
		endBlockHeight,
		parallelism,
	)

	return newNextBlockFunc(blockQueue, errorQueue)
}
