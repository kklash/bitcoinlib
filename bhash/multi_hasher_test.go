package bhash

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"

	"golang.org/x/crypto/ripemd160"
)

func TestMultiHasher(t *testing.T) {
	type Fixture struct {
		hashes    []hash.Hash
		inputHex  string
		outputHex string
	}

	fixtures := []Fixture{
		{
			[]hash.Hash{sha256.New(), sha256.New()},
			"deadbeef",
			"281dd50f6f56bc6e867fe73dd614a73c55a647a479704f64804b574cafb0f5c5",
		},
		{
			[]hash.Hash{sha256.New(), ripemd160.New()},
			"ffff",
			"e6abebacc6bf964f5131e80b241e3fe14bc3e156",
		},
		{
			[]hash.Hash{sha256.New()},
			"00000000",
			"df3f619804a92fdb4057192dc43dd748ea778adc52bc498ce80524c014b81119",
		},
		{
			[]hash.Hash{md5.New(), sha256.New()},
			"baddad",
			"cce9da9c59b7e16a8db3e09c1ac11ec475fb474892349bae4a8c0954e1cdf475",
		},
		{
			[]hash.Hash{sha512.New(), md5.New(), sha256.New(), sha512.New()},
			"baddad",
			"e7d1fc870552ae3727d0ccfb031c1528ee7ff5d5546280a428b43373057b5c8d39fbe10a4568652eac18228bd4805880a5d7fb2893bc50b1c6c79ef4517f2138",
		},
	}

	for _, fixture := range fixtures {
		mh := NewMultiHasher(fixture.hashes...)
		n, err := mh.Write(hex2bytes(fixture.inputHex))
		if err != nil {
			t.Errorf("Error writing input to MultiHasher: %s", fixture.inputHex)
			continue
		}
		if n != len(fixture.inputHex)/2 {
			t.Errorf("received unexpected number of bytes written\nWanted %d\nGot    %d", len(fixture.inputHex)/2, n)
			continue
		}

		output := mh.Sum(make([]byte, 0))
		if !bytes.Equal(output, hex2bytes(fixture.outputHex)) {
			t.Errorf("Received unexpected multihasher output\nWanted %s\nGot    %x", fixture.outputHex, output)
			continue
		}
	}
}
