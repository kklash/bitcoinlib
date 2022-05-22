package address

import (
	"testing"
)

func TestMakeP2SH(t *testing.T) {
	type Fixture struct {
		script  []byte
		address string
	}

	fixtures := []Fixture{
		Fixture{
			hex2bytes("0014ba397f72f8bd14fc26090a69d7f88c132d148d24"),
			"37HcHhcLJH9M1qnK3pkfqZiPYzUtNpGCGF",
		},
		Fixture{
			hex2bytes("5221037070dc9c046ce2cd18b5847899a016fafdd94f61c4032a91c8b33111c43f3bce21022bf784c5e4fee1dc7944fd723356304b4b307689fb9aa1a56c95177d6811d05321025257a485751eb4e9a264d9bc8632b9e5c25c24f6beffc159ce7b08e110ebe4d653ae"),
			"3FMk81WK7gChvBskTL7CDZZFnorr3s9Ju2",
		},
	}

	for _, fixture := range fixtures {
		addr := MakeP2SHFromScript(fixture.script)
		if addr != fixture.address {
			t.Errorf("P2SH address does not match fixture\nwanted %s\ngot %s", fixture.address, addr)
		}
	}
}
