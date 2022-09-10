package address

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/kklash/bitcoinlib/constants"
)

func hex2bytes(h string) []byte {
	data, _ := hex.DecodeString(h)
	return data
}

func TestMakeAndDecode(t *testing.T) {
	type Fixture struct {
		input        []byte
		hash         []byte
		format       constants.AddressFormat
		network      constants.Network
		address      string
		scriptPubKey []byte
	}

	fixtures := []Fixture{
		Fixture{
			input:        hex2bytes("030c5002a402c02c6662f7339f97150be7dfd8e266d8cf3374f0c823bc9c97ce53"),
			hash:         hex2bytes("7fd3519d11e8e4cfcaf55e3655af32161b24e023"),
			format:       constants.FormatP2PKH,
			network:      constants.BitcoinNetwork,
			address:      "1Cesy3vCXX2cf5TP3RN9dpJ7No8UXQJweg",
			scriptPubKey: hex2bytes("76a9147fd3519d11e8e4cfcaf55e3655af32161b24e02388ac"),
		},
		Fixture{
			input:        hex2bytes("0350123ebc957d3c117bf1d8119cf7432f15e86ebdbe027791a637e74e58259727"),
			hash:         hex2bytes("ce912eee3905de1ac5db35e891035e9f310c0e2a"),
			format:       constants.FormatP2WPKH,
			network:      constants.BitcoinNetwork,
			address:      "bc1qe6gjam3eqh0p43wmxh5fzq67nucscr32q35lpe",
			scriptPubKey: hex2bytes("0014ce912eee3905de1ac5db35e891035e9f310c0e2a"),
		},
		Fixture{
			input:        hex2bytes("0014262b2d2c13538ef75927cf7fd22d41103be19796"),
			hash:         hex2bytes("83197cecd6f3277be9ce00aea61b7b1f0ac34b87"),
			format:       constants.FormatP2SH,
			network:      constants.BitcoinNetwork,
			address:      "3DeCycLUMgbqs2MvJfdToCxp4tgREUa4iC",
			scriptPubKey: hex2bytes("a91483197cecd6f3277be9ce00aea61b7b1f0ac34b8787"),
		},
		Fixture{
			input:        hex2bytes("02217346dea363f34864a58fca85a9295ff6142622a0a78a5078b8b0dd9ffccf00"),
			hash:         hex2bytes("e79a6783b1546a18f7ef12c949d751b189cbf0e3"),
			format:       constants.FormatP2PKH,
			network:      constants.ZcashNetwork,
			address:      "t1ezD4rJqDBfnEQVMWqws9Ddi5f3g9jQiAA",
			scriptPubKey: hex2bytes("76a914e79a6783b1546a18f7ef12c949d751b189cbf0e388ac"),
		},
		Fixture{
			input:        hex2bytes("522103bb79cf5d6be36f6d598149a53d09c6fd1c8170241e3d63e5306b5b020ddb447a21022d1d6a538147bebcf0602defa6dd6373ee2e3d33be31f36ea7648ae20b096b7652ae"),
			hash:         hex2bytes("712c0646422a6e2c5bbb4f474f5c690d281167d627e855aeac29a33da1839191"),
			format:       constants.FormatP2WSH,
			network:      constants.BitcoinNetwork,
			address:      "bc1qwykqv3jz9fhzckamfar57hrfp55pze7kyl59tt4v9x3nmgvrjxgs5xwjnw",
			scriptPubKey: hex2bytes("0020712c0646422a6e2c5bbb4f474f5c690d281167d627e855aeac29a33da1839191"),
		},
	}

	for _, fixture := range fixtures {
		constants.CurrentNetwork = fixture.network

		encodedFromData, err := Make(fixture.format, fixture.input)
		if err != nil {
			t.Errorf("failed to make address %s - %s", fixture.address, err)
			continue
		}

		encodedFromHash, err := MakeFromHash(fixture.format, fixture.hash)
		if err != nil {
			t.Errorf("failed to make address %s - %s", fixture.address, err)
			continue
		}

		if encodedFromData != fixture.address {
			t.Errorf(
				"Make did not get expected result\nwanted %s\ngot %s",
				fixture.address,
				encodedFromData,
			)
			continue
		} else if encodedFromHash != fixture.address {
			t.Errorf(
				"MakeFromHash did not get expected result\nwanted %s\ngot %s",
				fixture.address,
				encodedFromHash,
			)
			continue
		}

		format, scriptPubKey, err := Decode(fixture.address)
		if err != nil {
			t.Errorf("Decode failed on '%s': %s", fixture.address, err)
			continue
		}

		if format != fixture.format {
			t.Errorf("Decode did not find correct address format; wanted '%s', got '%s'", fixture.format, format)
			continue
		}

		if !bytes.Equal(scriptPubKey, fixture.scriptPubKey) {
			t.Errorf(
				"Decode did not get expected result\nwanted %x\ngot    %x",
				fixture.scriptPubKey,
				scriptPubKey,
			)
			continue
		}
	}

	constants.CurrentNetwork = constants.BitcoinNetwork
}

func ExampleMake() {
	publicKey, _ := hex.DecodeString("03162c58483e649004a22bd8ed91ac88f172a0ccf640b34a3af7f470125e15b4bb")
	address, err := Make(constants.FormatP2WPKH, publicKey)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("addr: %s\n", address)
	// output:
	// addr: bc1qrkwwdpvkrq5etewn5k87lqgf3m3jjd5l9jdetv
}

func ExampleDecode() {
	format, scriptPubKey, err := Decode("bc1qrkwwdpvkrq5etewn5k87lqgf3m3jjd5l9jdetv")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s script pub key: %x\n", format, scriptPubKey)
	// output:
	// P2WPKH script pub key: 00141d9ce68596182995e5d3a58fef81098ee329369f
}

func TestDecodeBase58Address(t *testing.T) {
	type Fixture struct {
		address string
		version uint16
		payload []byte
	}

	fixtures := []Fixture{
		Fixture{
			address: "1Ckdc8grRFTh7hAtzcxMjAWMu6Kgf3Hb9W",
			version: 0,
			payload: hex2bytes("80e9d3b01aa679689c05a3bf96665e78b2db80c5"),
		},
		Fixture{
			address: "1J1j4XWyUytNDydKWLjEoBkhCsoe28ePUj",
			version: 0,
			payload: hex2bytes("ba9d3fb874c44ee404103a70600595c20af2cb9d"),
		},
		Fixture{
			address: "LaMmecUkVhooihC7gv8HagMUFuFR6k77hL",
			version: 0x30,
			payload: hex2bytes("a6045ac03781a1a083755586a3236ac18b096441"),
		},
		Fixture{
			address: "MKYx4DReDgpExkZuRZkyMi1eM2qv81FQio",
			version: 0x32,
			payload: hex2bytes("7fcf06ab7493a3900846c051218ff7d828e06f6a"),
		},
		Fixture{
			address: "t1Po9mYToY4Agp2VvQS255XeazUbLQcMWKt",
			version: 0x1cb8,
			payload: hex2bytes("40f9304afeba97d9045afc43640334429ea6d51b"),
		},
	}

	for _, fixture := range fixtures {
		version, hash, err := DecodeBase58Address(fixture.address)
		if err != nil {
			t.Errorf("Failed to decode address: %s", err)
			continue
		}

		if version != fixture.version {
			t.Errorf("decoded address version does not match\nwanted %d\ngot %d", fixture.version, version)
			continue
		}

		expectedHash := hex.EncodeToString(fixture.payload)
		actualHash := hex.EncodeToString(hash[:])

		if expectedHash != actualHash {
			t.Errorf("decoded address payload does not match\nwanted %s\ngot %s", expectedHash, actualHash)
			continue
		}
	}
}

func TestDecodeBech32Address(t *testing.T) {
	type Fixture struct {
		address string
		hrp     string
		version byte
		payload []byte
	}

	fixtures := []Fixture{
		{
			address: "bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej",
			hrp:     "bc",
			version: 0,
			payload: hex2bytes("701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d"),
		},
		{
			address: "bc1qwgczcvxwdea8qpm646jwgz9c2czuh63cgpw3k6",
			hrp:     "bc",
			version: 0,
			payload: hex2bytes("72302c30ce6e7a70077aaea4e408b85605cbea38"),
		},
		{
			address: "ltc1q5pg8p84y84af2vtpe0lcxamem3yqjnac9qdxrl",
			hrp:     "ltc",
			version: 0,
			payload: hex2bytes("a050709ea43d7a953161cbff837779dc48094fb8"),
		},
	}

	for _, fixture := range fixtures {
		hrp, version, payload, err := DecodeBech32Address(fixture.address)
		if err != nil {
			t.Errorf("failed to decode bech32 address: %s", err)
			continue
		}

		if hrp != fixture.hrp {
			t.Errorf("bech32 hrp did not match\nWanted %s\nGot    %s", fixture.hrp, hrp)
			continue
		}

		if version != fixture.version {
			t.Errorf("bech32 address version did not match\nWanted %d\nGot    %d", fixture.version, version)
			continue
		}

		if !bytes.Equal(payload, fixture.payload) {
			t.Errorf("bech32 address payload did not match\nWanted %x\nGot    %x", fixture.payload, payload)
			continue
		}
	}
}
