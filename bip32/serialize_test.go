package bip32

import (
	"encoding/hex"
	"testing"
)

type SerializationFixture struct {
	serializedKey     string
	key               string
	chainCode         string
	parentFingerprint string
	depth             byte
	index             uint32
	version           uint32
	private           bool
}

func (sf *SerializationFixture) Test(t *testing.T) {
	sf.TestSerialize(t)
	sf.TestDeserialize(t)
}

func (sf *SerializationFixture) TestSerialize(t *testing.T) {
	key, _ := hex.DecodeString(sf.key)
	chainCode, _ := hex.DecodeString(sf.chainCode)
	parentFingerprint, _ := hex.DecodeString(sf.parentFingerprint)

	serializedKey := serialize(key, chainCode, parentFingerprint, sf.depth, sf.index, sf.version, sf.private)
	if serializedKey != sf.serializedKey {
		t.Errorf("Expected serialized key %s\ngot %s", sf.serializedKey, serializedKey)
		return
	}
}

func (sf *SerializationFixture) TestDeserialize(t *testing.T) {
	key, chainCode, parentFingerprint, depth, index, version, err := Deserialize(sf.serializedKey)
	if err != nil {
		t.Errorf("Failed to deserialize key %s\nError: %s", sf.serializedKey, err)
		return
	}

	if actualKey := hex.EncodeToString(key); actualKey != sf.key {
		t.Errorf("Expected to deserialize key %s\ngot %s", sf.key, actualKey)
		return
	} else if actualChainCode := hex.EncodeToString(chainCode); actualChainCode != sf.chainCode {
		t.Errorf("Expected to deserialize chain code %s\ngot %s", sf.chainCode, actualChainCode)
		return
	} else if actualFingerprint := hex.EncodeToString(parentFingerprint); actualFingerprint != sf.parentFingerprint {
		t.Errorf("Expected to deserialize fingerprint %s\ngot %s", sf.parentFingerprint, actualFingerprint)
		return
	} else if depth != sf.depth {
		t.Errorf("Expected to deserialize depth %d\ngot %d", sf.depth, depth)
		return
	} else if index != sf.index {
		t.Errorf("Expected to deserialize index %d\ngot %d", sf.index, index)
		return
	} else if version != sf.version {
		t.Errorf("Expected to deserialize version %d\ngot %d", sf.version, version)
		return
	}
}

func TestSerializePublic(t *testing.T) {
	fixtures := []*SerializationFixture{
		{
			serializedKey:     "xpub6ASuArnXKPbfEwhqN6e3mwBcDTgzisQN1wXN9BJcM47sSikHjJf3UFHKkNAWbWMiGj7Wf5uMash7SyYq527Hqck2AxYysAA7xmALppuCkwQ",
			key:               "03501e454bf00751f24b1b489aa925215d66af2234e3891c3b21a52bedb3cd711c",
			chainCode:         "2a7857631386ba23dacac34180dd1983734e444fdbf774041578e9b6adb37c19",
			parentFingerprint: "5c1bd648",
			depth:             2,
			index:             1,
			version:           0x0488b21e,
			private:           false,
		},
		{
			serializedKey:     "xprv9uPDJpEQgRQfDcW7BkF7eTya6RPxXeJCqCJGHuCJ4GiRVLzkTXBAJMu2qaMWPrS7AANYqdq6vcBcBUdJCVVFceUvJFjaPdGZ2y9WACViL4L",
			key:               "491f7a2eebc7b57028e0d3faa0acda02e75c33b03c48fb288c41e2ea44e1daef",
			chainCode:         "e5fea12a97b927fc9dc3d2cb0d1ea1cf50aa5a1fdc1f933e8906bb38df3377bd",
			parentFingerprint: "41d63b50",
			depth:             1,
			index:             0x80000000,
			version:           0x0488ade4,
			private:           true,
		},
		{
			serializedKey:     "xpub6FnCn6nSzZAw5Tw7cgR9bi15UV96gLZhjDstkXXxvCLsUXBGXPdSnLFbdpq8p9HmGsApME5hQTZ3emM2rnY5agb9rXpVGyy3bdW6EEgAtqt",
			key:               "024d902e1a2fc7a8755ab5b694c575fce742c48d9ff192e63df5193e4c7afe1f9c",
			chainCode:         "9452b549be8cea3ecb7a84bec10dcfd94afe4d129ebfd3b3cb58eedf394ed271",
			parentFingerprint: "31a507b8",
			depth:             5,
			index:             2,
			version:           0x0488b21e,
			private:           false,
		},
	}

	for _, fixture := range fixtures {
		fixture.Test(t)
	}
}
