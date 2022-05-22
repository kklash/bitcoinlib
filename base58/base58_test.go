package base58

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestEncode(t *testing.T) {
	test := func(inputHex, expectedOutput string) {
		input, _ := hex.DecodeString(inputHex)
		output := Encode(input)
		if output != expectedOutput {
			t.Errorf("Failed to produce expected output\nWanted %s\nGot    %s", expectedOutput, output)
			return
		}

		decoded, err := Decode(output)
		if err != nil {
			t.Errorf("Failed to base58 decode %s\nError: %s", output, err)
			return
		}

		if !bytes.Equal(decoded, input) {
			t.Errorf("Failed to decode expected output\nWanted %x\nGot    %x", input, decoded)
			return
		}
	}

	fixtures := []struct{ input, output string }{
		{"61", "2g"},
		{"626262", "a3gV"},
		{"636363", "aPEr"},
		{"73696d706c792061206c6f6e6720737472696e67", "2cFupjhnEsSn59qHXstmK2ffpLv2"},
		{"00eb15231dfceb60925886b67d065299925915aeb172c06647", "1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L"},
		{"516b6fcd0f", "ABnLTmg"},
		{"bf4f89001e670274dd", "3SEo3LWLoPntC"},
		{"572e4794", "3EFU7m"},
		{"ecac89cad93923c02321", "EJDM8drfXA6uyA"},
		{"10c8511e", "Rt5zm"},
		{"00000000000000000000", "1111111111"},
		{"aeef6cae3dd0497efbe178c2d009c2211832c7cc615cb502df9cf59d3251", "c42xQenEFexu8H85Jia3cAhD7bJxECZmwZs8jqhCt"},
		{"", ""},
	}

	for _, fixture := range fixtures {
		test(fixture.input, fixture.output)
	}
}
