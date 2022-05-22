package base58check

import (
	"encoding/hex"
	"testing"
)

func hex2bytes(h string) []byte {
	buf, _ := hex.DecodeString(h)
	return buf
}

func TestEncode(t *testing.T) {
	type Fixture struct {
		input  []byte
		output string
	}

	fixtures := []Fixture{
		Fixture{
			input:  hex2bytes("009ff6dc18d42785f42ceaa72a4b757916ca7529e7"),
			output: "1FapAcp3HdoLop3QgVzoLcFAadijpRp6r3",
		},
		Fixture{
			input:  hex2bytes("001d17527311cc5454a181f70ebc686bb71a4951be"),
			output: "13epY1WbVbdMxEFsr9aQA74eyGcvPDE4SP",
		},
		Fixture{
			input:  hex2bytes("0573499684deccad704bb45d3d9ecfe5110de2b3b0"),
			output: "3CCbnavXL3qVfJcEjfcgoCsrVKPc4vM6qe",
		},
		Fixture{
			input:  hex2bytes("1cb8a6531e78151d273bd6dfeccf27f106d43dbc4d24"),
			output: "t1Z33kPandQXWpPGp5R7CVZfAdjRSbmiSZV",
		},
	}

	for _, fixture := range fixtures {
		encoded := Encode(fixture.input)
		if encoded != fixture.output {
			t.Errorf(
				"base58-check encoding did not output expected result\nwanted %s\ngot %s",
				fixture.output,
				encoded,
			)
		}
	}
}

func TestEncodeBase58CheckVersion(t *testing.T) {
	type Fixture struct {
		input   []byte
		version uint16
		output  string
	}

	fixtures := []Fixture{
		Fixture{
			input:   hex2bytes("9ff6dc18d42785f42ceaa72a4b757916ca7529e7"),
			version: 0,
			output:  "1FapAcp3HdoLop3QgVzoLcFAadijpRp6r3",
		},
		Fixture{
			input:   hex2bytes("1d17527311cc5454a181f70ebc686bb71a4951be"),
			version: 0,
			output:  "13epY1WbVbdMxEFsr9aQA74eyGcvPDE4SP",
		},
		Fixture{
			input:   hex2bytes("73499684deccad704bb45d3d9ecfe5110de2b3b0"),
			version: 5,
			output:  "3CCbnavXL3qVfJcEjfcgoCsrVKPc4vM6qe",
		},
		Fixture{
			input:   hex2bytes("a6531e78151d273bd6dfeccf27f106d43dbc4d24"),
			version: 0x1cb8,
			output:  "t1Z33kPandQXWpPGp5R7CVZfAdjRSbmiSZV",
		},
	}

	for _, fixture := range fixtures {
		encoded := EncodeVersion(fixture.input, fixture.version)
		if encoded != fixture.output {
			t.Errorf(
				"base58-check versioned encoding did not output expected result\nwanted %s\ngot %s",
				fixture.output,
				encoded,
			)
		}
	}
}
