package address

import (
	"testing"
)

func TestMakeP2WSHFromScript(t *testing.T) {
	type Fixture struct {
		script  []byte
		address string
	}

	fixtures := []Fixture{
		Fixture{
			script:  hex2bytes("0014ba397f72f8bd14fc26090a69d7f88c132d148d24"),
			address: "bc1qrm3a3hxkxvl7f4n6avxh68l3tm9v6npz2d25sat9c8zd7k2cegrsjykanx",
		},
		Fixture{
			script:  hex2bytes("5221037070dc9c046ce2cd18b5847899a016fafdd94f61c4032a91c8b33111c43f3bce21022bf784c5e4fee1dc7944fd723356304b4b307689fb9aa1a56c95177d6811d05321025257a485751eb4e9a264d9bc8632b9e5c25c24f6beffc159ce7b08e110ebe4d653ae"),
			address: "bc1qy6t2c8ldda0826ak5d4vjykejm7lvqkfwp2y9j9xr0eyplzc8peqaxd8rk",
		},
	}

	for _, fixture := range fixtures {
		addr, err := MakeP2WSHFromScript(fixture.script)
		if err != nil {
			t.Errorf(err.Error())
			continue
		}

		if addr != fixture.address {
			t.Errorf("P2WSH address does not match fixture\nwanted %s\ngot %s", fixture.address, addr)
		}
	}
}
