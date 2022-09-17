package bip39_test

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/kklash/bitcoinlib/bip39"
)

// Test vectors sourced from:
//
//	https://github.com/trezor/python-mnemonic/blob/e3b883007019c5156762aee8cdc3a31f3fc82c80/vectors.json
//
//go:embed test_vectors.json
var testVectorsJSON []byte

type Hex []byte

func (h *Hex) UnmarshalJSON(hexJSON []byte) error {
	var hexString string
	if err := json.Unmarshal(hexJSON, &hexString); err != nil {
		return err
	}
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return err
	}
	*h = Hex(decoded)
	return nil
}

type BIP39TestVector struct {
	Entropy  Hex
	Mnemonic string
	Seed     Hex
}

var testVectors []*BIP39TestVector

func init() {
	if err := json.Unmarshal(testVectorsJSON, &testVectors); err != nil {
		panic(fmt.Sprintf("failed to decode BIP39 test vectors: %s", err))
	}
}

// Confirms encoding of entropy works.
func TestEncodeToWords(t *testing.T) {
	for _, testVector := range testVectors {
		words, err := bip39.EncodeToWords(testVector.Entropy)
		if err != nil {
			t.Errorf("failed to encode entropy %x under BIP39: %s", testVector.Entropy, err)
			return
		}

		mnemonic := strings.Join(words, " ")
		if mnemonic != testVector.Mnemonic {
			t.Errorf("invalid mnemonic\nWanted %s\nGot    %s", testVector.Mnemonic, mnemonic)
			continue
		}
	}
}

// Confirms decoding of mnemonics works.
func TestDecodeWords(t *testing.T) {
	for _, testVector := range testVectors {
		words := strings.Split(testVector.Mnemonic, " ")
		entropy, err := bip39.DecodeWords(words)
		if err != nil {
			t.Errorf("failed to decode valid bip39 mnemonic: %s", err)
			return
		}

		if !bytes.Equal(entropy, testVector.Entropy) {
			t.Errorf("entropy does not match expected\nWanted %x\nGot    %x", testVector.Entropy, entropy)
			continue
		}
	}
}

// Confirms hashing of a mnemonic into a seed works.
func TestDeriveSeed(t *testing.T) {
	for _, testVector := range testVectors {
		words := strings.Split(testVector.Mnemonic, " ")
		seed := bip39.DeriveSeed(words, "TREZOR")

		if !bytes.Equal(seed, testVector.Seed) {
			t.Errorf("derived seed does not match expected\nWanted %x\nGot    %x", testVector.Seed, seed)
			continue
		}
	}
}

func BenchmarkDecodeWords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bip39.DecodeWords([]string{
			"scheme", "spot", "photo", "card", "baby", "mountain",
			"device", "kick", "cradle", "pact", "join", "borrow",
		})
	}
}

func BenchmarkEncodeToWords(b *testing.B) {
	entropy, _ := hex.DecodeString("c0ba5a8e914111210f2bd131f3d5e08d")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bip39.EncodeToWords(entropy)
	}
}
