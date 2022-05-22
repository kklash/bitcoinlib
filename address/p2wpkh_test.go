package address

import (
	"testing"
)

func TestMakeP2WPKHFromPublicKey(t *testing.T) {
	type Fixture struct {
		publicKey []byte
		address   string
	}

	fixtures := []Fixture{
		Fixture{
			hex2bytes("0365bc3ce47ae596a9c487762d0bc975ef746d0040a585495cd8454b4ed4ef6949"),
			"bc1qhguh7uhch520cfsfpf5a07yvzvk3frfy7vcmx6",
		},
	}

	for _, fixture := range fixtures {
		addr, err := MakeP2WPKHFromPublicKey(fixture.publicKey)
		if err != nil {
			t.Errorf(err.Error())
			continue
		}

		if addr != fixture.address {
			t.Errorf("P2WPKH address does not match fixture\nwanted %s\ngot %s", fixture.address, addr)
		}
	}
}
