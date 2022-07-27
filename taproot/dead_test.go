package taproot

import (
	"testing"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/ecc"
	"github.com/kklash/ekliptic"
)

func TestDeadGenerator(t *testing.T) {
	g := ecc.SerializePointUncompressed(ekliptic.Secp256k1_GeneratorX, ekliptic.Secp256k1_GeneratorY)
	hashed := bhash.Sha256(g)
	hx, hy, err := ecc.DeserializePoint(hashed[:])
	if err != nil {
		t.Errorf("failed to lift base point hash: %s", err)
		return
	}

	if hx.Cmp(deadGeneratorX) != 0 {
		t.Errorf(
			"failed to generate expected dead point X coordinate\nWanted %.64x\nGot    %.64x",
			deadGeneratorX, hx,
		)
	}
	if hy.Cmp(deadGeneratorY) != 0 {
		t.Errorf(
			"failed to generate expected dead point Y coordinate\nWanted %.64x\nGot    %.64x",
			deadGeneratorY, hy,
		)
	}
}
