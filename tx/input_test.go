package tx

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestInputEncodeDecode(t *testing.T) {
	type Fixture struct {
		raw       string // hex
		script    string // hex
		sequence  uint32
		prevHash  string // hex
		prevIndex uint32
		prevOut   string
	}

	fixtures := []Fixture{
		Fixture{
			raw:       "b85b41523f91890c7a665c473c1b5557332890c3576eca25fc60fa22e3c1451c050000001716001440af41ec0a8392051976b49f735827b04eda0381feffffff",
			script:    "16001440af41ec0a8392051976b49f735827b04eda0381",
			sequence:  0xfffffffe,
			prevHash:  "b85b41523f91890c7a665c473c1b5557332890c3576eca25fc60fa22e3c1451c",
			prevIndex: 5,
			prevOut:   "1c45c1e322fa60fc25ca6e57c390283357551b3c475c667a0c89913f52415bb8:5",
		},
		Fixture{
			raw:       "cb6fa8149033c17979144d55602e591f3f8ac24d742bc2716569867afbcc9d260000000017160014633c1568569ec6e5f9a0034936fb85675f5b8945feffffff",
			script:    "160014633c1568569ec6e5f9a0034936fb85675f5b8945",
			sequence:  0xfffffffe,
			prevHash:  "cb6fa8149033c17979144d55602e591f3f8ac24d742bc2716569867afbcc9d26",
			prevIndex: 0,
			prevOut:   "269dccfb7a86696571c22b744dc28a3f1f592e60554d147979c1339014a86fcb:0",
		},
		Fixture{
			raw:       "cb6fa8149033c17979144d55602e591f3f8ac24d742bc2716569867afbcc9d260000000000ffffffff",
			script:    "",
			sequence:  0xffffffff,
			prevHash:  "cb6fa8149033c17979144d55602e591f3f8ac24d742bc2716569867afbcc9d26",
			prevIndex: 0,
			prevOut:   "269dccfb7a86696571c22b744dc28a3f1f592e60554d147979c1339014a86fcb:0",
		},
	}

	for _, fixture := range fixtures {
		data, _ := hex.DecodeString(fixture.raw)
		input, err := InputFromReader(bytes.NewReader(data))
		if err != nil {
			t.Errorf(err.Error())
			continue
		}
		scriptHex := hex.EncodeToString(input.Script)
		prevHashHex := hex.EncodeToString(input.PrevOut.Hash[:])

		if scriptHex != fixture.script {
			t.Errorf("script does not match\nwanted %s\ngot    %s", fixture.script, scriptHex)
		} else if prevHashHex != fixture.prevHash {
			t.Errorf("prev hash hex does not match\nwanted %s\ngot    %s", fixture.prevHash, prevHashHex)
		} else if input.Sequence != fixture.sequence {
			t.Errorf("sequence does not match\nwanted %d\ngot    %d", fixture.sequence, input.Sequence)
		} else if input.PrevOut.Index != fixture.prevIndex {
			t.Errorf("prev out index does not match\nwanted %d\ngot    %d", fixture.prevIndex, input.PrevOut.Index)
		} else if input.PrevOut.String() != fixture.prevOut {
			t.Errorf("prev out string does not match\nwanted %s\ngot    %s", fixture.prevOut, input.PrevOut)
		}

		bytesRead, err := input.WriteTo(new(bytes.Buffer))
		if err != nil {
			t.Errorf("Failed to encode input: %s", err)
			continue
		} else if expectedBytesRead := int64(len(fixture.raw) / 2); bytesRead != expectedBytesRead {
			t.Errorf("bytes read did not match\nwanted %d\ngot %d", expectedBytesRead, bytesRead)
			continue
		}

		encodedInputHex := hex.EncodeToString(input.Bytes())
		if encodedInputHex != fixture.raw {
			t.Errorf("serialized input does not match fixture\nwanted %s\ngot %s", fixture.raw, encodedInputHex)
			continue
		}

		clone := input.Clone()
		if !reflect.DeepEqual(input, clone) {
			t.Errorf("input.Clone() did not produce exact copy")
			continue
		}

		size := input.Size()
		if size != len(data) {
			t.Errorf("encoded size did not match expected\nWanted %d\nGot    %d", len(data), size)
			continue
		}
	}
}
