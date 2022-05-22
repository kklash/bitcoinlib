package bip32

import (
	"encoding/hex"
	"testing"

	"github.com/kklash/bitcoinlib/ecc"
)

type KeyCompressionFixture struct {
	compressedHex   string
	uncompressedHex string
}

func (kf *KeyCompressionFixture) TestCompress(t *testing.T) {
	uncompressed, _ := hex.DecodeString(kf.uncompressedHex)
	if err := ValidatePublicKeyBytes(uncompressed); err != nil {
		t.Errorf("Expected public key %x to be valid", uncompressed)
		return
	}

	compressed, err := CompressPublicKeyBytes(uncompressed)
	if err != nil {
		t.Errorf("Failed to compress public key %x\nError: %s", uncompressed, err)
		return
	} else if hex.EncodeToString(compressed) != kf.compressedHex {
		t.Errorf("Expected compressed public key %s\ngot %x", kf.compressedHex, compressed)
		return
	}
}

func (kf *KeyCompressionFixture) TestUncompress(t *testing.T) {
	compressed, _ := hex.DecodeString(kf.compressedHex)
	if !IsCompressedPublicKey(compressed) {
		t.Errorf("Expected public key %x to be interpreted as compressed", compressed)
		return
	} else if err := ValidatePublicKeyBytes(compressed); err != nil {
		t.Errorf("Expected public key %x to be valid", compressed)
		return
	}

	x, y, err := UncompressPublicKey(compressed)
	if err != nil {
		t.Errorf("Failed to uncompress public key %x\nError: %s", compressed, err)
		return
	} else if uncompressedHex := hex.EncodeToString(ecc.SerializePoint(x, y)); uncompressedHex != kf.uncompressedHex {
		t.Errorf("Expected uncompressed public key %s\ngot %s", kf.uncompressedHex, uncompressedHex)
		return
	}
}

func TestPublicKeyCompression(t *testing.T) {
	fixtures := []*KeyCompressionFixture{
		{
			"020000000000000000000000000000000000000000000000000000000000000001",
			"00000000000000000000000000000000000000000000000000000000000000014218f20ae6c646b363db68605822fb14264ca8d2587fdd6fbc750d587e76a7ee",
		},
		{
			"0254236f7d1124fc07600ad3eec5ac47393bf963fbf0608bcce255e685580d16d9",
			"54236f7d1124fc07600ad3eec5ac47393bf963fbf0608bcce255e685580d16d92b266380164d9b76f0f442d0a618314ae40a463e0ca69f68d40d572c21f8927e",
		},
		{
			"02588d202afcc1ee4ab5254c7847ec25b9a135bbda0f2bc69ee1a714749fd77dc9",
			"588d202afcc1ee4ab5254c7847ec25b9a135bbda0f2bc69ee1a714749fd77dc9f88ff2a00d7e752d44cbe16e1ebcf0890b76ec7c78886109dee76ccfc8445424",
		},
		{
			"02098190f2263db5049e659448343683a77efcfe5162b445dc491dcf19c34ecccd",
			"098190f2263db5049e659448343683a77efcfe5162b445dc491dcf19c34ecccda31663c11d630c9c4c5cf14997ea7df9f774082829fd373ca9be4a4739e44ba8",
		},
		{
			"022f5379552d5f68b83747057a901ca7b8846337683236e6055ed2ce5f72230503",
			"2f5379552d5f68b83747057a901ca7b8846337683236e6055ed2ce5f7223050325ac711e204d5ef4349f2eb18ba03fdafb7b9a830159ca2c05c08c3359dcfff4",
		},
	}

	for _, fixture := range fixtures {
		fixture.TestCompress(t)
		fixture.TestUncompress(t)
	}
}
