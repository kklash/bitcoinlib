package bech32

import (
	"testing"
)

func TestEncodeBech32(t *testing.T) {
	type Fixture struct {
		input   []byte
		version byte
		prefix  string
		output  string
	}

	fixtures := []Fixture{
		Fixture{
			input:   hex2bytes("ce912eee3905de1ac5db35e891035e9f310c0e2a"),
			version: 0,
			prefix:  "bc",
			output:  "bc1qe6gjam3eqh0p43wmxh5fzq67nucscr32q35lpe",
		},
		Fixture{
			input:   hex2bytes("38fd8bb6c6d6e38215312ec5d145f1bdee902788"),
			version: 0,
			prefix:  "bc",
			output:  "bc1q8r7chdkx6m3cy9f39mzaz303hhhfqfug8jxdud",
		},
		Fixture{
			input:   hex2bytes("98455a7c8b006f586f8eafeeb6c9a096ed601c84"),
			version: 0,
			prefix:  "dgb",
			output:  "dgb1qnpz45lytqph4smuw4lhtdjdqjmkkq8yyk3m0h3",
		},
	}

	for _, fixture := range fixtures {
		encoded, err := Encode(fixture.prefix, fixture.version, fixture.input)
		if err != nil {
			t.Errorf("Failed to encode bech32 string: %x\nError: %s", fixture.input, err)
			continue
		}

		if encoded != fixture.output {
			t.Errorf(
				"bech32 encoding did not output expected result\nwanted %s\ngot    %s",
				fixture.output,
				encoded,
			)
			continue
		}
	}
}
