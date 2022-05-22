package tx

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestOutputEncodeDecode(t *testing.T) {
	type Fixture struct {
		raw    string // hex
		script string // hex
		value  uint64
	}

	fixtures := []Fixture{
		Fixture{
			raw:    "40380f000000000017a914cc6e1f626e5670871663d95475dd102f1c9a8a7587", // P2SH
			script: "a914cc6e1f626e5670871663d95475dd102f1c9a8a7587",
			value:  0xf3840,
		},
		Fixture{
			raw:    "d67f0400000000001976a9141aae47f952b6d2415ff898e476f13f45c3963aa388ac", // P2PKH
			script: "76a9141aae47f952b6d2415ff898e476f13f45c3963aa388ac",
			value:  0x47fd6,
		},
		Fixture{
			raw:    "0000000000000000116a0f6d79207377656174792062616c6c73", // OP_RETURN
			script: "6a0f6d79207377656174792062616c6c73",
			value:  0,
		},
		Fixture{
			raw:    "47dd50000000000016001415047c1bfe4409edae828d50280367b1a099a87b",
			script: "001415047c1bfe4409edae828d50280367b1a099a87b",
			value:  5299527,
		},
	}

	for _, fixture := range fixtures {
		data, _ := hex.DecodeString(fixture.raw)
		output, err := OutputFromReader(bytes.NewReader(data))
		if err != nil {
			t.Errorf(err.Error())
			continue
		}

		scriptHex := hex.EncodeToString(output.Script)
		if scriptHex != fixture.script {
			t.Errorf("script does not match\nwanted %s\ngot %s", fixture.script, scriptHex)
			continue
		} else if output.Value != fixture.value {
			t.Errorf("value does not match\nwanted %d\ngot %d", fixture.value, output.Value)
			continue
		}

		bytesRead, err := output.WriteTo(new(bytes.Buffer))
		if err != nil {
			t.Errorf("Failed to serialize: %s", err)
			continue
		} else if int(bytesRead) != len(fixture.raw)/2 {
			t.Errorf("bytes read count from WriteTo does not match\nwanted %d\ngot %d", len(fixture.raw)/2, bytesRead)
			continue
		}

		encodedOutputHex := hex.EncodeToString(output.Bytes())

		if encodedOutputHex != fixture.raw {
			t.Errorf("serialized output does not match fixture\nwanted %s\ngot %s", fixture.raw, encodedOutputHex)
			continue
		}

		clone := output.Clone()
		if !reflect.DeepEqual(output, clone) {
			t.Errorf("output.Clone() did not produce exact copy")
			continue
		}

		clone.Script[0] = 0xff
		if output.Script[0] == clone.Script[0] {
			t.Errorf("clone.Script still points to original output.Script")
			continue
		}

		size := output.Size()
		if size != len(data) {
			t.Errorf("output size does not match\nWanted %d\nGot    %d", len(data), size)
			continue
		}
	}
}
