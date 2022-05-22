package tx

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestWitnessEncodeDecode(t *testing.T) {
	type Fixture struct {
		raw     string   // hex
		witness []string // hex
	}

	fixtures := []Fixture{
		Fixture{
			raw: "0247304402206b077561ec75a66ec5969bc81db1ac831bf36fff281323fe1861be4630bb29dc022049b885ee89bd69e37c2deab62132a602dac4a6c15a1683e088560375ebdf6d8b012103776852b68dc55cb66019d44e0f12323dd42329cfafc5e3af6409cb4ef2e2acd7",
			witness: []string{
				"304402206b077561ec75a66ec5969bc81db1ac831bf36fff281323fe1861be4630bb29dc022049b885ee89bd69e37c2deab62132a602dac4a6c15a1683e088560375ebdf6d8b01",
				"03776852b68dc55cb66019d44e0f12323dd42329cfafc5e3af6409cb4ef2e2acd7",
			},
		},
		Fixture{
			raw:     "00",
			witness: []string{},
		},
		Fixture{
			raw: "040047304402206c420b92a4b2b635b7b5bb05e7c7be9b5b44c59f2e311f19ff65a716cfe69278022035ebd8106515a15945e9442f84ed28f3fbfba355336c9423184ec0272e6015b201473044022049d252a30101696b75973ffd011e08c88ab6cbbd82e5871ecb0c9e844d71eff402203a72cfec7173256206882e3c8424726727f8774925d6abee00ac952fb3ef26ba016952210257f116736c36e79ee83e9d561bb68a7b40bcdd5a3302d5ac8bb7692ce97e4d312103889a3b32dfa72227948a606efc1d7c85e72fcfdcfb019cf251db87796730870321035f67179191cf55360fefab85753b4a53db0b1b4ed494262d8f9efd0ee92c729453ae",
			witness: []string{
				"",
				"304402206c420b92a4b2b635b7b5bb05e7c7be9b5b44c59f2e311f19ff65a716cfe69278022035ebd8106515a15945e9442f84ed28f3fbfba355336c9423184ec0272e6015b201",
				"3044022049d252a30101696b75973ffd011e08c88ab6cbbd82e5871ecb0c9e844d71eff402203a72cfec7173256206882e3c8424726727f8774925d6abee00ac952fb3ef26ba01",
				"52210257f116736c36e79ee83e9d561bb68a7b40bcdd5a3302d5ac8bb7692ce97e4d312103889a3b32dfa72227948a606efc1d7c85e72fcfdcfb019cf251db87796730870321035f67179191cf55360fefab85753b4a53db0b1b4ed494262d8f9efd0ee92c729453ae",
			},
		},
	}

	for _, fixture := range fixtures {
		data, _ := hex.DecodeString(fixture.raw)
		witness, err := WitnessFromReader(bytes.NewReader(data))
		if err != nil {
			t.Errorf(err.Error())
			continue
		}

		if len(witness) != len(fixture.witness) {
			t.Errorf("witness chunk length does not match\nwanted %d\ngot %d", len(fixture.witness), len(witness))
			continue
		}

		for i := 0; i < len(fixture.witness); i++ {
			chunkHex := hex.EncodeToString(witness[i])
			if chunkHex != fixture.witness[i] {
				t.Errorf("witness chunk does not match\nwanted %s\ngot %s", fixture.witness[i], chunkHex)
				continue
			}
		}

		bytesRead, err := witness.WriteTo(new(bytes.Buffer))
		if err != nil {
			t.Errorf("Failed to encode witness: %s", err)
			continue
		} else if expectedBytesRead := int64(len(fixture.raw) / 2); bytesRead != expectedBytesRead {
			t.Errorf("bytes read did not match expected\nwanted %d\ngot %d", expectedBytesRead, bytesRead)
			continue
		}

		encodedWitnessHex := hex.EncodeToString(witness.Bytes())
		if encodedWitnessHex != fixture.raw {
			t.Errorf("serialized witness does not match fixture\nwanted %s\ngot %s", fixture.raw, encodedWitnessHex)
			continue
		}

		clone := witness.Clone()
		if !reflect.DeepEqual(witness, clone) {
			t.Errorf("witness.Clone() did not produce exact copy")
			continue
		}

		if len(witness) > 1 {
			clone[1][0] = 0xff
			if witness[1][0] == clone[1][0] {
				t.Errorf("clone still points to original witness")
				continue
			}
		}
	}
}
