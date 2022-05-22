package address

import (
	"testing"
)

func TestMakeP2PKH(t *testing.T) {
	type Fixture struct {
		publicKey []byte
		address   string
	}

	fixtures := []Fixture{
		Fixture{
			publicKey: hex2bytes("021dc77bed63ee01e7c52dc3a7fcf77d7b669d89e9f7f901b173e92c29f868b8cd"),
			address:   "1MrVmxzEgVn6rJSaRXyJsk4ho2Auve7yWx",
		},
		Fixture{
			publicKey: hex2bytes("02536c6c25f5e97da80d1a1d962ebb0150d761d4eb024948917dd28c7846cb8ac4"),
			address:   "1GpxKcGAeajtHRAMb3HBp6VMoqTKAnw6tU",
		},
		Fixture{
			publicKey: hex2bytes("02ef814d9311852a1c9074b37cbaf50a7bc77310c9f217bbb37d54db119f0b2ff5"),
			address:   "1KcKwYfg3JPchCLCUCHMLkmuo2i4tXGorn",
		},
	}

	for _, fixture := range fixtures {
		addr, err := MakeP2PKHFromPublicKey(fixture.publicKey)
		if err != nil {
			t.Errorf(err.Error())
			continue
		}

		if addr != fixture.address {
			t.Errorf("P2PKH address does not match fixture\nwanted %s\ngot %s", fixture.address, addr)
		}
	}
}
