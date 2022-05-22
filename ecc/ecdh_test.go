package ecc

import (
	"bytes"
	"testing"

	"github.com/kklash/ekliptic"
)

func TestSharedSecret(t *testing.T) {
	curve := new(ekliptic.Curve)

	fixtures := []struct {
		key1 string
		key2 string
	}{
		{
			"70ACD14CFEA8509F41584A3C166B6D7ABCE52A10BA860BCFB26129E06CC00086",
			"6AB5319B31B8612029557F154BE94ABA8DCF80AA04C55DC4537F30BC02370798",
		},
		{
			"2E0FC5B257C3F9680CC5881FADA8D5367213864A1906E272E1C7242D35C915DF",
			"334CC7A5E871E5863E7D337E63B3A7BB69993CBADA1F6970954480C30E7D7822",
		},
		{
			"592534B42EBB6565EE891EE513A15FB2CF3530369205133D981809603BB68311",
			"C375B2D0500ED24E60B383748E17B78BAE6DC8F80C998888518A2A3D055FC1CE",
		},
		{
			"0000000000000000000000000000000000000000000000000000000000000001",
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364140",
		},
	}

	for _, fixture := range fixtures {
		priv1, priv2 := pointAtHex(fixture.key1, fixture.key2)

		pubX1, pubY1 := curve.ScalarBaseMult(priv1.Bytes())
		pubX2, pubY2 := curve.ScalarBaseMult(priv2.Bytes())

		secret1 := SharedSecret(priv1, pubX2, pubY2)
		secret2 := SharedSecret(priv2, pubX1, pubY1)

		if !bytes.Equal(secret1, secret2) {
			t.Errorf("Failed to form matching secrets\nGot\n %x\n %x", secret1, secret2)
			return
		}

		if len(secret1) != 32 {
			t.Errorf("wrong shared secret size: %d", len(secret1))
		}
	}
}
