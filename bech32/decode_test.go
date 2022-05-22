package bech32

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func hex2bytes(h string) []byte {
	d, _ := hex.DecodeString(h)
	return d
}

func TestDecode(t *testing.T) {
	type Fixture struct {
		input   string
		prefix  string
		version byte
		output  []byte
	}

	fixtures := []Fixture{
		Fixture{
			input:   "bc1qe6gjam3eqh0p43wmxh5fzq67nucscr32q35lpe",
			prefix:  "bc",
			version: 0,
			output:  hex2bytes("ce912eee3905de1ac5db35e891035e9f310c0e2a"),
		},
		Fixture{
			input:   "bc1q8r7chdkx6m3cy9f39mzaz303hhhfqfug8jxdud",
			prefix:  "bc",
			version: 0,
			output:  hex2bytes("38fd8bb6c6d6e38215312ec5d145f1bdee902788"),
		},
		Fixture{
			input:   "dgb1qnpz45lytqph4smuw4lhtdjdqjmkkq8yyk3m0h3",
			prefix:  "dgb",
			version: 0,
			output:  hex2bytes("98455a7c8b006f586f8eafeeb6c9a096ed601c84"),
		},
	}

	for _, fixture := range fixtures {
		prefix, version, decoded, err := Decode(fixture.input)
		if err != nil {
			t.Errorf("failed to decode bech32 string: %s", err)
			continue
		}

		if prefix != fixture.prefix {
			t.Errorf("bech32 prefix did not match expected\nwanted %s\ngot    %s", fixture.prefix, prefix)
			continue
		}

		if version != fixture.version {
			t.Errorf("bech32 version did not match expected\nwanted %d\ngot    %d", fixture.version, version)
			continue
		}

		if !bytes.Equal(decoded, fixture.output) {
			t.Errorf(
				"bech32 decoding did not output expected result\nwanted %x\ngot    %x",
				fixture.output,
				decoded,
			)
			continue
		}
	}
}
