package bhash

import (
	"crypto/sha256"
)

// NewTaggedHasher returns a function which generates a hash whose preimage is prepended twice
// with the hash of the given tag string.
//
//  tagged_hash(tag, ...chunks) = sha256(sha256(tag) || sha256(tag) || chunks...)
func NewTaggedHasher(tag string) func(...[]byte) []byte {
	hashedTag := sha256.Sum256([]byte(tag))
	return func(chunks ...[]byte) []byte {
		h := sha256.New()
		h.Write(hashedTag[:])
		h.Write(hashedTag[:])
		for _, chunk := range chunks {
			h.Write(chunk)
		}
		return h.Sum(nil)
	}
}
