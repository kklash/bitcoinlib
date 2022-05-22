package bip38

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestEncrypt(t *testing.T) {
	type Fixture struct {
		password     string
		encryptedKey string
		privateKey   string
		compressed   bool
		ecMultiply   bool
	}

	fixtures := []Fixture{
		{
			password:     "TestingOneTwoThree",
			encryptedKey: "6PYNKZ1EAgYgmQfmNVamxyXVWHzK5s6DGhwP4J5o44cvXdoY7sRzhtpUeo",
			privateKey:   "CBF4B9F70470856BB4F40F80B87EDB90865997FFEE6DF315AB166D713AF433A5",
			compressed:   true,
			ecMultiply:   false,
		},
	}

	for _, fixture := range fixtures {
		privateKey, _ := hex.DecodeString(fixture.privateKey)

		encryptedKey, err := encrypt(
			privateKey,
			fixture.password,
			fixture.compressed,
			fixture.ecMultiply,
			false,
		)

		if err != nil {
			t.Errorf("Failed to encrypt key: %s", err)
			return
		}

		if encryptedKey != fixture.encryptedKey {
			t.Errorf("encrypted key did not match expected\nWanted %s\nGot    %s", fixture.encryptedKey, encryptedKey)
			return
		}

		decryptedPrivateKey, compressed, err := decrypt(encryptedKey, fixture.password)
		if err != nil {
			t.Errorf("failed to decrypt key: %s", err)
			return
		}

		if !bytes.Equal(decryptedPrivateKey, privateKey) {
			t.Errorf("failed to decrypt expected private key\nWanted %x\nGot    %x", privateKey, decryptedPrivateKey)
			return
		}

		if compressed != fixture.compressed {
			t.Errorf("failed to parse expected compressed flag\nWanted %v\nGot    %v", fixture.compressed, compressed)
			return
		}
	}
}
