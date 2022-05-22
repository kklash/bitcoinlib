package blockscan

import "github.com/kklash/bitcoinlib/blocks"

// NextBlockFunc is a function which is returned from the streaming methods of BlockScanner.
// It returns the next block in the queue, or an error if one was encountered. When no more
// blocks are available, it returns two nil values.
type NextBlockFunc func() (*blocks.Block, error)

func newNextBlockFunc(blockQueue chan *blocks.Block, errorQueue chan error) NextBlockFunc {
	return func() (*blocks.Block, error) {
		select {
		case block, more := <-blockQueue:
			if !more {
				return nil, nil
			}

			return block, nil

		case err := <-errorQueue:
			return nil, err
		}
	}
}
