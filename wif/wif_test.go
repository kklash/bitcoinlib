package wif

import (
	"encoding/hex"
	"testing"
)

type Fixture struct {
	privkey    string
	version    byte
	compressed bool
	wif        string
	valid      bool
}

var fixtures = []Fixture{
	Fixture{
		privkey:    "58e63ffbd1f8ee5c9cb15d0c1087fb53a31172ff7731cfd4c5deb34935388b5f",
		version:    0x80,
		compressed: true,
		wif:        "KzCX6cDLTVW21eng76dcSKfJemm7q2rpZc838VRkkNx7R7AswR36",
		valid:      true,
	},
	Fixture{
		privkey:    "74f4ec234609204c4aef8a0c1800f2f79e769fc3d833c111e8a0ebb1cd89a3f2",
		version:    0x80,
		compressed: true,
		wif:        "L194QLTCZhHJxWiEBJwdVQ5CuJBJuriUzYq4jQszuT2GGp9VmggM",
		valid:      true,
	},
	Fixture{
		privkey:    "b8613ae0d0ac185f0d7100322e927e5fff59e1050b30f0af3fb51dbb698e7a9e",
		version:    0x80,
		compressed: false,
		wif:        "5KDVKHoNT64tDi1EFDWp6pEv7GgqYY2i63ha44frKhiFFjLPe2s",
		valid:      true,
	},
	Fixture{
		privkey:    "c5add55a1a6ab37b2599c5753049841b51bd5eb589c392a1e8a76ed92005560a",
		version:    0xef,
		compressed: false,
		wif:        "935ycGYCaWLCgAv6tydzPp1FvK3sbpFZXguFoMfRSSfQCtCnUvq",
		valid:      true,
	},
	Fixture{
		privkey:    "",
		version:    0,
		compressed: false,
		wif:        "abeE9dFj2d938d",
		valid:      false,
	},
}

func TestEncode(t *testing.T) {
	for _, fixture := range fixtures {
		if !fixture.valid {
			continue
		}

		privkey, _ := hex.DecodeString(fixture.privkey)

		var encoded string
		var err error
		if fixture.compressed {
			encoded, err = Encode(privkey, fixture.version)
		} else {
			encoded, err = EncodeUncompressed(privkey, fixture.version)
		}

		if err != nil {
			t.Errorf("failed to encode WIF key: %s", err)
			continue
		}

		if encoded != fixture.wif {
			t.Errorf(
				"WIF encoding did not output expected result\nwanted %s\ngot %s",
				fixture.wif,
				encoded,
			)
		}
	}
}

func TestDecode(t *testing.T) {
	for _, fixture := range fixtures {
		privkey, version, compressed, err := Decode(fixture.wif)

		if !fixture.valid {
			if err == nil {
				t.Errorf("failed to return error when decoding invalid string")
			}
			continue
		}

		if err != nil {
			t.Errorf("failed to decode WIF key: %s", err)
			continue
		}

		if version != fixture.version {
			t.Errorf("WIF decoded version byte does not match\nwanted %d\ngot %d", fixture.version, version)
		}

		if compressed != fixture.compressed {
			t.Errorf("WIF decoded compression flag does not match\nwanted %t\ngot %t", fixture.compressed, compressed)
		}

		if hex.EncodeToString(privkey) != fixture.privkey {
			t.Errorf(
				"WIF decoded key did not output expected result\nwanted %s\ngot %x",
				fixture.privkey,
				privkey,
			)
		}
	}
}

func TestValidate(t *testing.T) {
	for _, fixture := range fixtures {
		valid := Validate(fixture.wif)
		if valid != fixture.valid {
			t.Errorf("WIF validation failed for key %s\nwanted %t\ngot %t", fixture.wif, fixture.valid, valid)
		}
	}
}
