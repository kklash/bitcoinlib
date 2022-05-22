package der

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/kklash/bitcoinlib/constants"
)

func bigIntFromHex(h string) *big.Int {
	data, _ := hex.DecodeString(h)
	return new(big.Int).SetBytes(data)
}

func TestEncodeDecodeSignature(t *testing.T) {
	fixtures, err := readFixtures()
	if err != nil {
		t.Errorf("failed to read fixtures: %s", err)
		return
	}

	for _, fixture := range fixtures {
		r := bigIntFromHex(fixture.R)
		s := bigIntFromHex(fixture.S)

		encoded, err := EncodeSignature(r, s, constants.SigHashAll)
		if err != nil {
			t.Errorf("Failed to encode signature: %s", err)
			return
		}

		expected, _ := hex.DecodeString(fixture.DEREncoded)

		if !bytes.Equal(encoded, expected) {
			t.Errorf("did not receive expected DER-encoded signature\nWanted %x\nGot    %x", expected, encoded)
			return
		}

		decodedR, decodedS, sigHashType, err := DecodeSignature(encoded)
		if err != nil {
			t.Errorf("failed to decode DER signature: %s", err)
			return
		}

		if sigHashType != constants.SigHashAll {
			t.Errorf("unexpected sighash from decoded signature\nWanted %d\nGot    %d", constants.SigHashAll, sigHashType)
			return
		}

		if decodedR.Cmp(r) != 0 {
			t.Errorf("failed to parse signature r value\nWanted %x\nGot    %x", r, decodedR)
			return
		}

		if decodedS.Cmp(s) != 0 {
			t.Errorf("failed to parse signature r value\nWanted %x\nGot    %x", s, decodedS)
			return
		}
	}
}
