package script

import (
	"bytes"
	"testing"
)

func TestMakeP2MS(t *testing.T) {
	type Fixture struct {
		publicKeys   [][]byte
		keysRequired uint32
		scriptPubKey []byte
	}

	fixtures := []Fixture{
		Fixture{
			publicKeys: [][]byte{
				hex2bytes("030264a09adddcd9d3e139807809d441b3a60fec1dd6029c2770dd16de9dca2ef9"),
				hex2bytes("03d801595232f3c5384b186b7309d207716153d94f8603eac18134f11586bbc54d"),
			},
			keysRequired: 2,
			scriptPubKey: hex2bytes("5221030264a09adddcd9d3e139807809d441b3a60fec1dd6029c2770dd16de9dca2ef92103d801595232f3c5384b186b7309d207716153d94f8603eac18134f11586bbc54d52ae"),
		},
		Fixture{
			publicKeys: [][]byte{
				hex2bytes("030264a09adddcd9d3e139807809d441b3a60fec1dd6029c2770dd16de9dca2ef9"),
				hex2bytes("03d801595232f3c5384b186b7309d207716153d94f8603eac18134f11586bbc54d"),
				hex2bytes("030c7b4d3b504d91cec06a7c81a3a254317d9bd1f18ba03cd69b41a5c9b8907b50"),
				hex2bytes("03e8dbdab8eb77ba938247ef780baaff4db5d44968f66eb72f74a8be3b36544eaa"),
				hex2bytes("034fd37b3f0fb4c166074db52d9d93406b92dd18e0be3452c6dfebc86db6496b3b"),
				hex2bytes("02f2989d776babe11df783d58f5c0f7994b083446d296bc37be2f0c318668b4054"),
				hex2bytes("02adc11d60f68ad4c29edddfcb0917266723be3875454f19ee50313bda8370d002"),
				hex2bytes("025d24befc600fe9738ecf18c2c43d993c539e091e78a2eaf6da43b71d382eb5fc"),
				hex2bytes("02011b7c0e74bc0e25d457ec8cbc16ed8cc36d673f37099f431a2c8519c0f6b50c"),
				hex2bytes("037928186ca1c68f78349aa641b677dabcdf05d34da8117d55edda722deb52a798"),
				hex2bytes("02b2ba4368bb4716049cb1e7bce1d9a3ab6dd712ab64c70e473bb9422e9cc1f506"),
				hex2bytes("02fb8e6cd5bc3ffa99bc0b4d1580390ee9f3f9496862a42ed46151d7af3f699032"),
				hex2bytes("02b6bd06cab132e1b18dabf64d4e015325c3d0229d262491aa45f69549dc32df37"),
				hex2bytes("03f4dcb77dfbf065581fcaa13c3c6845334d0ce74c297273072e3a9b666a09eba6"),
				hex2bytes("03b1222b065f84dec69e26aa3febcec6ed3dc69be9c2dab52e16960f50a289a259"),
				hex2bytes("02107af12a8e474ccfee349362b2820fd380f0f909b3734d7ee9cc89acfcab616d"),
				hex2bytes("0304b4832c370ff6ef21db7e9ebd43516a1b9ea935c60c312a55bcf27d830f9fee"),
				hex2bytes("0377c8baaaea9825606970fca2a4abfb940bfdcc6040b0212243fb9d63c185dccb"),
			},
			keysRequired: 3,
			scriptPubKey: hex2bytes("5321030264a09adddcd9d3e139807809d441b3a60fec1dd6029c2770dd16de9dca2ef92103d801595232f3c5384b186b7309d207716153d94f8603eac18134f11586bbc54d21030c7b4d3b504d91cec06a7c81a3a254317d9bd1f18ba03cd69b41a5c9b8907b502103e8dbdab8eb77ba938247ef780baaff4db5d44968f66eb72f74a8be3b36544eaa21034fd37b3f0fb4c166074db52d9d93406b92dd18e0be3452c6dfebc86db6496b3b2102f2989d776babe11df783d58f5c0f7994b083446d296bc37be2f0c318668b40542102adc11d60f68ad4c29edddfcb0917266723be3875454f19ee50313bda8370d00221025d24befc600fe9738ecf18c2c43d993c539e091e78a2eaf6da43b71d382eb5fc2102011b7c0e74bc0e25d457ec8cbc16ed8cc36d673f37099f431a2c8519c0f6b50c21037928186ca1c68f78349aa641b677dabcdf05d34da8117d55edda722deb52a7982102b2ba4368bb4716049cb1e7bce1d9a3ab6dd712ab64c70e473bb9422e9cc1f5062102fb8e6cd5bc3ffa99bc0b4d1580390ee9f3f9496862a42ed46151d7af3f6990322102b6bd06cab132e1b18dabf64d4e015325c3d0229d262491aa45f69549dc32df372103f4dcb77dfbf065581fcaa13c3c6845334d0ce74c297273072e3a9b666a09eba62103b1222b065f84dec69e26aa3febcec6ed3dc69be9c2dab52e16960f50a289a2592102107af12a8e474ccfee349362b2820fd380f0f909b3734d7ee9cc89acfcab616d210304b4832c370ff6ef21db7e9ebd43516a1b9ea935c60c312a55bcf27d830f9fee210377c8baaaea9825606970fca2a4abfb940bfdcc6040b0212243fb9d63c185dccb0112ae"),
		},
	}

	for _, fixture := range fixtures {
		scriptPubKey := MakeP2MS(fixture.keysRequired, fixture.publicKeys...)
		if !bytes.Equal(scriptPubKey, fixture.scriptPubKey) {
			t.Errorf("script pub key does not match:\n wanted %x\n got %x", fixture.scriptPubKey, scriptPubKey)
		}
	}
}

func TestRedeemP2MS(t *testing.T) {
	signature := hex2bytes("3046022100a07b2821f96658c938fa9c68950af0e69f3b2ce5f8258b3a6ad254d4bc73e11e022100e82fab8df3f7e7a28e91b3609f91e8ebf663af3a4dc2fd2abd954301a5da67e701")

	expectedScriptSig := append(hex2bytes("0049"), signature...)

	actualScriptSig := RedeemP2MS(signature)

	if !bytes.Equal(actualScriptSig, expectedScriptSig) {
		t.Errorf("P2MS redeem script did not build correctly\nWanted %x\nGot    %x", expectedScriptSig, actualScriptSig)
	}
}
