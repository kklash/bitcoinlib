package base58check

import (
	"bytes"
	"testing"
)

func TestDecodeBase58Check(t *testing.T) {
	type Fixture struct {
		input  string
		output []byte
	}

	fixtures := []Fixture{
		Fixture{
			input:  "1FapAcp3HdoLop3QgVzoLcFAadijpRp6r3",
			output: hex2bytes("009ff6dc18d42785f42ceaa72a4b757916ca7529e7"),
		},
		Fixture{
			input:  "13epY1WbVbdMxEFsr9aQA74eyGcvPDE4SP",
			output: hex2bytes("001d17527311cc5454a181f70ebc686bb71a4951be"),
		},
		Fixture{
			input:  "3CCbnavXL3qVfJcEjfcgoCsrVKPc4vM6qe",
			output: hex2bytes("0573499684deccad704bb45d3d9ecfe5110de2b3b0"),
		},
		Fixture{
			input:  "L3RzENM63ELkfRoXxYaBzcvkGYXBhCmGJvyN1hTT7g234tDxAZym",
			output: hex2bytes("80b95707c7d8db7006204f6bef366987de804f295a09ca0121435eac03dd226fc101"),
		},
		Fixture{
			input:  "t1Z33kPandQXWpPGp5R7CVZfAdjRSbmiSZV",
			output: hex2bytes("1cb8a6531e78151d273bd6dfeccf27f106d43dbc4d24"),
		},
	}

	for _, fixture := range fixtures {
		decoded, err := Decode(fixture.input)
		if err != nil {
			t.Errorf("failed to decode base58-check string: %s", err)
			continue
		}

		if !bytes.Equal(decoded, fixture.output) {
			t.Errorf(
				"base58-check decoding did not output expected result\nwanted %x\ngot %x",
				fixture.output,
				decoded,
			)
		}
	}
}
