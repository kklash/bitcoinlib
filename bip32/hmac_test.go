package bip32

import (
	"encoding/hex"
	"testing"
)

func TestHmacSha512(t *testing.T) {
	test := func(keyHex, dataHex, expectedHex string) {
		key, _ := hex.DecodeString(keyHex)
		data, _ := hex.DecodeString(dataHex)
		hash := hmacSha512(key, data)

		if hashHex := hex.EncodeToString(hash); hashHex != expectedHex {
			t.Errorf("HMAC-SHA512 test failed\n wanted %s\n got %s", expectedHex, hashHex)
		}
	}

	test(
		"0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b",
		"4869205468657265",
		"87aa7cdea5ef619d4ff0b4241a1d6cb02379f4e2ce4ec2787ad0b30545e17cdedaa833b7d6b8a702038b274eaea3f4e4be9d914eeb61f1702e696c203a126854",
	)

	test(
		"4a656665",
		"7768617420646f2079612077616e7420666f72206e6f7468696e673f",
		"164b7a7bfcf819e2e395fbe73b56e0a387bd64222e831fd610270cd7ea2505549758bf75c05a994a6d034f65f8f0e6fdcaeab1a34d4a6b4b636e070a38bce737",
	)

	test(
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"fa73b0089d56a284efb0f0756c890be9b1b5dbdd8ee81a3655f83e33b2279d39bf3e848279a722c806b485a47e67c807b946a337bee8942674278859e13292fb",
	)
}
