// Package blockscan provides a BlockScanner struct which can be used to perform high-throughput scanning
// of blocks from a bitcoind RPC connection. blockscan handles parsing blocks from the node into structs.
package blockscan

import (
	"encoding/hex"
	"strings"

	"github.com/kklash/bitcoinlib/blocks"
	"github.com/kklash/bitcoinlib/rpc"
)

// BlockScanner is used to fetch and stream blocks from an RPC connection. It can store a context which
// allows the caller to cancel block streaming.
type BlockScanner struct {
	Connection *rpc.Connection
}

// NewBlockScanner returns a pointer to a BlockScanner which will fetch blocks from the given rpc.Connection.
func NewBlockScanner(conn *rpc.Connection) *BlockScanner {
	return &BlockScanner{Connection: conn}
}

// GetBlockByHeight returns a pointer to a blocks.Block containing the block data at the given height.
func (scanner *BlockScanner) GetBlockByHeight(height uint32) (*blocks.Block, error) {
	blockHash, err := scanner.Connection.Request("getblockhash", height)
	if err != nil {
		return nil, err
	}

	blockHex, err := scanner.Connection.Request("getblock", blockHash, 0)
	if err != nil {
		return nil, err
	}

	hexReader := hex.NewDecoder(strings.NewReader(blockHex.(string)))

	block, err := blocks.FromReader(hexReader)
	if err != nil {
		return nil, err
	}

	return block, nil
}
