// Package bhash exposes commonly used utility functions for hashing data in Bitcoin.
package bhash

import (
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

// Sha256 returns the SHA256 hash of the given data.
func Sha256(data []byte) [32]byte {
	return sha256.Sum256(data)
}

// DoubleSha256 returns the double-SHA256 hash of the given data.
func DoubleSha256(data []byte) [32]byte {
	firstpass := Sha256(data)
	return Sha256(firstpass[:])
}

// Ripemd160 returns the rmd160 hash of the given data.
func Ripemd160(data []byte) (result [20]byte) {
	hash := ripemd160.New()
	hash.Write(data)
	copy(result[:], hash.Sum([]byte{}))
	return
}

// Hash160 returns the rmd160 hash of a SHA256 hash of the given data.
func Hash160(data []byte) [20]byte {
	firstPass := Sha256(data)
	return Ripemd160(firstPass[:])
}
