package ecc

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"os"
	"testing"
)

func TestSchnorr(t *testing.T) {
	fh, err := os.Open("schnorr_fixtures.csv")
	if err != nil {
		t.Errorf("failed to open schnorr test fixtures file: %s", err)
		return
	}
	defer fh.Close()
	rows, err := csv.NewReader(fh).ReadAll()
	if err != nil {
		t.Errorf("failed to read CSV rows: %s", err)
		return
	}
	rows = rows[1:]

	for i, columns := range rows {
		privateKey, err0 := hex.DecodeString(columns[1])
		publicKey, err1 := hex.DecodeString(columns[2])
		auxRand, err2 := hex.DecodeString(columns[3])
		messageHash, err3 := hex.DecodeString(columns[4])
		fixtureSig, err4 := hex.DecodeString(columns[5])
		shouldValidate := columns[6] == "TRUE"
		comment := columns[7]

		for _, err := range []error{err0, err1, err2, err3, err4} {
			if err != nil {
				t.Errorf("failed to decode schnorr fixture index %d (%q): %s", i, comment, err)
				continue
			}
		}

		if len(privateKey) > 0 {
			signature := SignSchnorr(privateKey, messageHash, auxRand)

			if !bytes.Equal(signature, fixtureSig) {
				t.Errorf("incorrect schnorr signature for %q\nWanted %x\nGot    %x", comment, fixtureSig, signature)
				continue
			}
		}

		valid := VerifySchnorr(publicKey, messageHash, fixtureSig)
		if valid != shouldValidate {
			t.Errorf("schnorr signature verification mismatch for %q; wanted %v, got %v", comment, shouldValidate, valid)
			return
		}
	}
}
