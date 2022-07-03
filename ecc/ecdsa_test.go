package ecc

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

func TestECDSA(t *testing.T) {
	var fixtures []map[string]string

	fixtureData, err := os.ReadFile("ecdsa_fixtures.json")
	if err != nil {
		t.Errorf("failed to read ECDSA fixtures: %s", err)
		return
	}

	if err := json.Unmarshal(fixtureData, &fixtures); err != nil {
		t.Errorf("failed to parse ECDSA fixtures: %s", err)
		return
	}

	for _, fixture := range fixtures {
		hash, _ := hex.DecodeString(fixture["hash"])
		privateKey, _ := hex.DecodeString(fixture["d"])
		expectedR, expectedS := pointAtHex(fixture["r"], fixture["s"])

		r, s := SignECDSA(privateKey, hash)
		if !equal(r, expectedR) || !equal(s, expectedS) {
			t.Errorf(
				"Expected (r, s) point:\n(\n %X,\n %X\n)\nGot point \n(\n %X,\n %X\n)",
				expectedR, expectedS,
				r, s,
			)
			return
		}
		publicKey := GetPublicKeyCompressed(privateKey)
		if !VerifyECDSA(hash, r, s, publicKey) {
			t.Errorf("Expected signature to be verified as valid:\n(\n %X,\n %X\n)", r, s)
		}

		publicKey = GetPublicKeyCompressed([]byte{7})
		if VerifyECDSA(hash, r, s, publicKey) {
			t.Errorf("Expected signature to be invalid with wrong pub key:\n(\n %X,\n %X\n)", r, s)
		}
	}
}
