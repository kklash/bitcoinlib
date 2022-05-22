// Package merkle provides a merkle-tree hashing implementation used for Bitcoin applications.
package merkle

import (
	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/common"
)

// Concatenate 32-byte hashes together. Used for joining hashes in merkle-tree rows.
func concatHashes(hashes ...[32]byte) []byte {
	joined := make([]byte, len(hashes)*32)
	for i := 0; i < len(hashes); i++ {
		copy(joined[i*32:], hashes[i][:])
	}
	return joined
}

// MerkleRootHashInternal generates the merkle root of a slice of transaction hashes
// (internal byte order), without performing any byte-order manipulation.
func MerkleRootHashInternal(hashes [][32]byte) [32]byte {
	if len(hashes) == 1 {
		// Done
		return hashes[0]
	} else if len(hashes) == 2 {
		return bhash.DoubleSha256(concatHashes(hashes...))
	}

	// Handle odd number of hashes
	if len(hashes)%2 == 1 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	// Recursively create the intermediate rows
	nextRow := make([][32]byte, len(hashes)/2)
	for i := 0; i < len(nextRow); i++ {
		k := i * 2
		joined := concatHashes(hashes[k], hashes[k+1])
		nextRow[i] = bhash.DoubleSha256(joined)
	}

	return MerkleRootHashInternal(nextRow)
}

// MerkleRootHash generates the merkle root for a slice of transaction IDs in RPC byte order,
// and returns the merkle root, also in RPC byte order. The TXIDs are reversed before being hashed using
// MerkleRootHashInternal, whose result is also reversed before being returned by this function.
func MerkleRootHash(txids [][32]byte) [32]byte {
	hashes := make([][32]byte, len(txids))
	for i, txid := range txids {
		copy(hashes[i][:], common.ReverseBytes(txid[:]))
	}

	merkle := MerkleRootHashInternal(hashes)
	common.ReverseBytesInPlace(merkle[:])
	return merkle
}
