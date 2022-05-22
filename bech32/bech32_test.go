package bech32

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestBech32(t *testing.T) {
	test := func(version byte, hrp, inputHex, expectedOutput string) {
		data, _ := hex.DecodeString(inputHex)

		encoded, err := Encode(hrp, version, data)

		if err != nil {
			t.Errorf("Failed to encode data: %s\nError: %s", inputHex, err)
			return
		}

		if encoded != expectedOutput {
			t.Errorf("Bech32 encoded output did not match\nWanted %s\nGot    %s", expectedOutput, encoded)
			return
		}

		decodedHrp, decodedVersion, decoded, err := Decode(encoded)
		if err != nil {
			t.Errorf("Failed to decode bech32 string: %s\nError: %s", encoded, err)
			return
		} else if decodedHrp != hrp {
			t.Errorf("Bech32 decoded HRP does not match\nWanted %s\nGot    %s", hrp, decodedHrp)
			return
		} else if !bytes.Equal(decoded, data) {
			t.Errorf("Bech32 decoded bytes do not match\nWanted %x\nGot    %x", data, decoded)
			return
		} else if decodedVersion != version {
			t.Errorf("Bech32 decoded version byte does not match\nWanted %d\nGot    %d", version, decodedVersion)
			return
		}
	}

	test(0, "bc", "751e76e8199196d454941c45d1b3a323f1433bd6", "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4")
	test(0, "bc", "ce912eee3905de1ac5db35e891035e9f310c0e2a", "bc1qe6gjam3eqh0p43wmxh5fzq67nucscr32q35lpe")
	test(0, "dgb", "98455a7c8b006f586f8eafeeb6c9a096ed601c84", "dgb1qnpz45lytqph4smuw4lhtdjdqjmkkq8yyk3m0h3")
	test(0, "tb", "1863143c14c5166804bd19203356da136c985678cd4d27a1b8c6329604903262", "tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sl5k7")
	test(
		1,
		"bc",
		"751e76e8199196d454941c45d1b3a323f1433bd6751e76e8199196d454941c45d1b3a323f1433bd6",
		"bc1pw508d6qejxtdg4y5r3zarvary0c5xw7kw508d6qejxtdg4y5r3zarvary0c5xw7k7grplx",
	)
	test(2, "bc", "751e76e8199196d454941c45d1b3a323", "bc1zw508d6qejxtdg4y5r3zarvaryvg6kdaj")
	test(16, "bc", "751e", "bc1sw50qa3jx3s")
	test(0, "bc", "000000000000", "bc1qqqqqqqqqqq576m3x")
	test(0, "bc", "0000000000", "bc1qqqqqqqqqqv9qus")
	test(0, "bc", "000000", "bc1qqqqqqdjyd6q")
	test(0, "bc", "000001", "bc1qqqqqz7h9wfd")
	test(
		0,
		"1",
		"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqsdxzd2",
	)
	test(
		0,
		"tb",
		"000000c4a5cad46221b2a187905e5266362b99d5e91c6ce24d165dab93e86433",
		"tb1qqqqqp399et2xygdj5xreqhjjvcmzhxw4aywxecjdzew6hylgvsesrxh6hy",
	)
}
