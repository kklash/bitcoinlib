package bhash

import (
	"hash"
)

// MultiHasher is a chain of cryptographic hashes, one feeding into the next from
// first to last. Think of it as an io.MultiWriter but for hash functions.
type MultiHasher struct {
	hashes []hash.Hash
}

// NewMultiHasher returns a pointer to a new MultiHasher, which will use the
// given hashes provided to hash data.
func NewMultiHasher(hashes ...hash.Hash) *MultiHasher {
	if len(hashes) == 0 {
		return nil
	}

	return &MultiHasher{hashes}
}

// BlockSize returns the block size of the first hash in the chain.
func (mh *MultiHasher) BlockSize() int {
	return mh.hashes[0].BlockSize()
}

// Size returns the output size of the last hash in the chain.
func (mh *MultiHasher) Size() int {
	return mh.hashes[len(mh.hashes)-1].Size()
}

// Reset resets the MultiHasher to its initial state.
func (mh *MultiHasher) Reset() {
	mh.hashes[0].Reset()
}

// Write writes the given data to the MultiHasher.
func (mh *MultiHasher) Write(p []byte) (int, error) {
	return mh.hashes[0].Write(p)
}

// Sum computes the final hash by piping output from each hash to the next in the chain,
// returning the result of appending the final hash to b. If b is nil, it simply returns
// the finalized hash.
func (mh *MultiHasher) Sum(b []byte) []byte {
	var hashed []byte
	for _, h := range mh.hashes {
		h.Write(hashed)
		hashed = h.Sum(nil)
	}
	return hashed
}
