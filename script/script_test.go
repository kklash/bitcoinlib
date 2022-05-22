package script

import (
	"bytes"
	"testing"

	"github.com/kklash/bitcoinlib/constants"
)

func TestStripOpCode(t *testing.T) {
	type Fixture struct {
		input  []byte
		output []byte
		op     byte
	}

	fixtures := []Fixture{
		{
			hex2bytes("01ab00ab"),
			hex2bytes("01ab00"),
			0xab,
		},
	}

	for _, fixture := range fixtures {
		actual, err := StripOpCode(fixture.input, fixture.op)
		if err != nil {
			t.Errorf("Failed to strip op code: %s", err)
			return
		}

		if !bytes.Equal(actual, fixture.output) {
			t.Errorf("failed to strip op code\nWanted %x\nGot    %x", fixture.output, actual)
			return
		}
	}
}

func TestDecompile(t *testing.T) {
	// TODO
}

func TestClassifyOutput(t *testing.T) {
	fixtures := []struct {
		script []byte
		format constants.AddressFormat
	}{
		{
			hex2bytes("00202696ac1fed6f5e756bb6a36ac912d996fdf602c9705442c8a61bf240fc583872"),
			constants.FormatP2WSH,
		},
		{
			hex2bytes("0014ce6a28589e056b0bdd67c464033677a7ac35ce05"),
			constants.FormatP2WPKH,
		},
		{
			hex2bytes("a9146505c5df8a4dcde5d6ef68cb4c15efa19c01991c87"),
			constants.FormatP2SH,
		},
		{
			hex2bytes("76a91491c79c05a31adead59033ebf47acab299b4cdba488ac"),
			constants.FormatP2PKH,
		},
		{
			hex2bytes("0014ce6a28589e056b0bdd67c464033677a7ac35ce0511"),
			constants.FormatNONSTANDARD,
		},
		{
			hex2bytes("5221030264a09adddcd9d3e139807809d441b3a60fec1dd6029c2770dd16de9dca2ef92103d801595232f3c5384b186b7309d207716153d94f8603eac18134f11586bbc54d52ae"), // P2MS
			constants.FormatNONSTANDARD,
		},
	}

	for _, fixture := range fixtures {
		format := ClassifyOutput(fixture.script)

		if format != fixture.format {
			t.Errorf("classified output script incorrectly\nWanted %s\nGot    %s", fixture.format, format)
			continue
		}
	}
}
