package tx

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"testing"
)

var (
	// Example segregated witness transaction.
	witnessTx = "02000000000103b85b41523f91890c7a665c473c1b5557332890c3576eca25fc60fa22e3c1451c050000001716001440af41ec0a8392051976b49f735827b04eda0381feffffff2a885bc4d4ced1c929bfa5bf39cab1f7147cf71cb11ac8ae4ac8180c57a7532700000000171600148848580b83c5abb8645643495e84e5310b781d19feffffffd9f98190a2145cfddf01ac254c6b69f6ec773f0562f48fedfaf5906cb74a9d9c0000000017160014135acfc3259d66830a0a12dfab4f1f7a358a67a2feffffff0106cb0b00000000001976a9145a55f9da2176d6325c4df2376b1cf1a82ba1f03488ac024730440220559914eb16f9103c842ef4e040d790d044a4843df305445bd770b64f129727c9022071da33beb433588d71502e6aa74e647a36e86209ade2bad909f233abec3a5e96012102d6c56b9a931c89f5b504436d40580858568391eef455169bc5c245c65b04b59d02473044022073671a3865432575c0926d081ee1e973e83c7d87ef543f1608978e5ea4e0250f02203ed619949c1513ba854c758612cfa2e68671ff0599bf2b9a7f5b9cca33ef3310012102a707207c5012b17a848feb1ba8ef54fda6d6cc55d1cba9344d811b36ca643ab602473044022074c73a3a75d36f205674af07add15933e7add57c45fa5d762cb7fb9f749926090220250c9176fb12bf2a37b5970616e0888d49d84924a4974311e35c259f55f57922012103fb10291d505f07cbe8655a896682235a9748f14e4be788eba2a140774768380234c80800"
	/*
		02000000 // version
		0001 // witness flag
		03 // nInputs

		b85b41523f91890c7a665c473c1b5557332890c3576eca25fc60fa22e3c1451c 05000000 // prevout 0
		17 // scriptlen 0
		16001440af41ec0a8392051976b49f735827b04eda0381 // script 0
		feffffff // sequence 0

		2a885bc4d4ced1c929bfa5bf39cab1f7147cf71cb11ac8ae4ac8180c57a75327 00000000 // prevout 1
		17 // scriptlen 1
		1600148848580b83c5abb8645643495e84e5310b781d19 // script 1
		feffffff // sequence 1

		d9f98190a2145cfddf01ac254c6b69f6ec773f0562f48fedfaf5906cb74a9d9c 00000000 // prevout 2
		17 // scriptlen 2
		160014135acfc3259d66830a0a12dfab4f1f7a358a67a2 // script 2
		feffffff // sequence 2


		01 // nOutputs

		06cb0b0000000000 // value 0
		19 // scriptPubKey length 0
		76a9145a55f9da2176d6325c4df2376b1cf1a82ba1f03488ac // scriptPubKey 0

		// witnesses

		02 // nWitnesses (input 0)
		47 // witness length 0
		30440220559914eb16f9103c842ef4e040d790d044a4843df305445bd770b64f129727c9022071da33beb433588d71502e6aa74e647a36e86209ade2bad909f233abec3a5e9601
		21 // witness length 1
		02d6c56b9a931c89f5b504436d40580858568391eef455169bc5c245c65b04b59d


		02 // nWitnesses (input 1)
		47 // witness length 0
		3044022073671a3865432575c0926d081ee1e973e83c7d87ef543f1608978e5ea4e0250f02203ed619949c1513ba854c758612cfa2e68671ff0599bf2b9a7f5b9cca33ef331001
		21 // witness length 1
		02a707207c5012b17a848feb1ba8ef54fda6d6cc55d1cba9344d811b36ca643ab6

		02 // nWitnesses 2 (input 2)
		47 // witness length 0
		3044022074c73a3a75d36f205674af07add15933e7add57c45fa5d762cb7fb9f749926090220250c9176fb12bf2a37b5970616e0888d49d84924a4974311e35c259f55f5792201
		21 // witness length 1
		03fb10291d505f07cbe8655a896682235a9748f14e4be788eba2a1407747683802

		34c80800 // locktime
	*/

	// Example legacy non-segwit transaction.
	legacyTx = "0200000001ba87e36602b41a1fdd16ec9956f620eed6de69979fe076d167de476ae695cde9000000006a47304402207542770b195c8e27a66f9072e353f361d428e9c691692d1fcf59cf245ff9a17302206c30d9e1b2e84ce7969b8912e5dd9e34da071b5f44561195c103ff8965bb7680012103939e84ab4d070debe46835c043fb5cf85d243dc9337cde76d5ae062b21dc0ccfffffffff01b8cf0700000000001976a914f498aef7c6f23b7718ea709cf7ccb3b84aae217788ac00000000"
	/*
		02000000 // version
		01 // nInputs

		ba87e36602b41a1fdd16ec9956f620eed6de69979fe076d167de476ae695cde9 00000000 // prevout 0
		6a // scriptlen 0
		47304402207542770b195c8e27a66f9072e353f361d428e9c691692d1fcf59cf245ff9a17302206c30d9e1b2e84ce7969b8912e5dd9e34da071b5f44561195c103ff8965bb7680012103939e84ab4d070debe46835c043fb5cf85d243dc9337cde76d5ae062b21dc0cc // script 0
		fffffffff // sequence 0

		01 // nOutputs

		b8cf070000000000 // value 0
		19 // scriptPubKey len 0
		76a914f498aef7c6f23b7718ea709cf7ccb3b84aae217788ac // scriptPubKey 0

		00000000 // locktime
	*/
)

func TestTxEncodeDecode(t *testing.T) {
	// only test overall structure layout. Tests for decoding specific
	// data structures (inputs, outputs, witnesses) is handled in sub-packages.
	type Fixture struct {
		raw         string
		txid        string
		version     int32
		nInputs     int
		nOutputs    int
		nWitnesses  int
		locktime    uint32
		weightUnits int
		vSize       int
	}

	fixtures := []Fixture{
		Fixture{
			raw:         witnessTx,
			txid:        "b666e22a593be35688d99115f6eca71c1b1737cf422d9f9b4cff904e00d0dde7",
			version:     2,
			nInputs:     3,
			nOutputs:    1,
			nWitnesses:  3,
			locktime:    575540,
			weightUnits: 1267,
			vSize:       317,
		},
		Fixture{
			raw:         legacyTx,
			txid:        "25184d820c0bad8a2e0d5811a0a8c9becdaadc5ec400c9942c2c056350671ed4",
			version:     2,
			nInputs:     1,
			nOutputs:    1,
			nWitnesses:  0,
			locktime:    0,
			weightUnits: 764,
			vSize:       191,
		},
	}

	for _, fixture := range fixtures {
		data, _ := hex.DecodeString(fixture.raw)
		tx, err := FromReader(bytes.NewReader(data))
		if err != nil {
			t.Errorf("failed to decode transaction: %s", err)
			continue
		}

		if tx.Version != fixture.version {
			t.Errorf("invalid TX version\nwanted %d\ngot %d", fixture.version, tx.Version)
			continue
		} else if tx.Locktime != fixture.locktime {
			t.Errorf("invalid TX locktime\nwanted %d\ngot %d", fixture.locktime, tx.Locktime)
			continue
		} else if len(tx.Inputs) != fixture.nInputs {
			t.Errorf("invalid number of inputs\nwanted %d\ngot %d", fixture.nInputs, len(tx.Inputs))
			continue
		} else if len(tx.Outputs) != fixture.nOutputs {
			t.Errorf("invalid number of outputs\nwanted %d\ngot %d", fixture.nOutputs, len(tx.Outputs))
			continue
		}

		if fixture.nWitnesses == 0 {
			if tx.Witnesses != nil {
				t.Errorf("expected nil witness")
				continue
			}
		} else {
			if tx.Witnesses == nil {
				t.Errorf("expected witness data - got nil")
				continue
			} else if len(tx.Witnesses) != fixture.nWitnesses {
				t.Errorf("invalid number of witnesses\nwanted %d\ngot %d", fixture.nWitnesses, len(tx.Witnesses))
				continue
			}
		}

		bytesRead, err := tx.WriteTo(new(bytes.Buffer))
		if err != nil {
			t.Errorf("Failed to encode tx: %s", err)
			continue
		} else if expectedBytesRead := int64(len(fixture.raw) / 2); bytesRead != expectedBytesRead {
			t.Errorf("bytes read did not match\nwanted %d\ngot %d", expectedBytesRead, bytesRead)
			continue
		}

		serializedTx := tx.Bytes()
		encodedTxHex := hex.EncodeToString(serializedTx)

		if encodedTxHex != fixture.raw {
			t.Errorf("serialized tx did not match fixture\nwanted %s\ngot %s", fixture.raw, encodedTxHex)
			continue
		}

		size := tx.Size()
		if size != len(serializedTx) {
			t.Errorf("tx size does not match expected\nwanted %d\ngot    %d", len(serializedTx), size)
			continue
		}

		if weightUnits := tx.WeightUnits(); weightUnits != fixture.weightUnits {
			t.Errorf("tx weightUnits does not match expected\nwanted %d\ngot    %d", fixture.weightUnits, weightUnits)
			return
		}

		if vSize := tx.VSize(); vSize != fixture.vSize {
			t.Errorf("tx vSize does not match expected\nwanted %d\ngot    %d", fixture.vSize, vSize)
			return
		}

		txid, err := tx.Id(false)
		if err != nil {
			t.Errorf("failed to get transaction id: %s", err)
			continue
		}

		if txid != fixture.txid {
			t.Errorf("txid did not match fixture\nwanted %s\ngot %s", fixture.txid, txid)
			continue
		}

		clone := tx.Clone()
		if !reflect.DeepEqual(tx, clone) {
			t.Errorf("tx.Clone() did not produce exact copy")
			fmt.Println(tx)
			fmt.Println()
			fmt.Println(clone)
			continue
		}
	}
}

func ExampleFromBytes() {
	txBytes, _ := hex.DecodeString("02000000000101d6e3bde671ba50371cc04a694d2fd637c25a732989bc242fa3d2f159998a3eee0000000017160014d6d26f092c608c6667d1a43d8fe13acc0d348e02ffffffff02a58ffa000000000017a91463bb8e6658d347563c6264c69ff81bbed7db59fb87a33c45000000000017a914807c70ab304dcad7437dcfa701b70db4f7981f7e870248304502210092ecb520c563e9530d3c2b0011e08d2c0b42f3a07cfa8d39e7edc457afcb558102203f1bc3f38ef4d542ba467b2edd7c379ee105e4adfb16bc616fb0a7ba8cd9f6a301210205aa55376f5c847e949ae077362a67d9d2061ecb92a038c508ddcd9b4608bb3e00000000")

	tx, err := FromBytes(txBytes)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Version: %d\n", tx.Version)
	fmt.Printf("Locktime: %d\n", tx.Locktime)
	fmt.Println("Inputs[0]:")
	fmt.Printf("  script: %x\n", tx.Inputs[0].Script)
	fmt.Printf("  sequence: %d\n", tx.Inputs[0].Sequence)
	fmt.Printf("  prevout: %x:%d\n", tx.Inputs[0].PrevOut.Hash, tx.Inputs[0].PrevOut.Index)
	fmt.Println("Outputs[0]:")
	fmt.Printf("  script: %x\n", tx.Outputs[0].Script)
	fmt.Printf("  value: %d\n", tx.Outputs[0].Value)

	// Output:
	// Version: 2
	// Locktime: 0
	// Inputs[0]:
	//   script: 160014d6d26f092c608c6667d1a43d8fe13acc0d348e02
	//   sequence: 4294967295
	//   prevout: d6e3bde671ba50371cc04a694d2fd637c25a732989bc242fa3d2f159998a3eee:0
	// Outputs[0]:
	//   script: a91463bb8e6658d347563c6264c69ff81bbed7db59fb87
	//   value: 16420773
}

func TestVSize(t *testing.T) {
	tx, err := FromBytes(mustHex("0100000000010115e180dc28a2327e687facc33f10f2a20da717e5548406f7ae8b4c811072f85603000000171600141d7cd6c75c2e86f4cbf98eaed221b30bd9a0b928ffffffff019caef505000000001976a9141d7cd6c75c2e86f4cbf98eaed221b30bd9a0b92888ac02483045022100f764287d3e99b1474da9bec7f7ed236d6c81e793b20c4b5aa1f3051b9a7daa63022016a198031d5554dbb855bdbe8534776a4be6958bd8d530dc001c32b828f6f0ab0121038262a6c6cec93c2d3ecd6c6072efea86d02ff8e3328bbd0242b20af3425990ac00000000"))
	if err != nil {
		t.Errorf("failed to decode transaction: %s", err)
		return
	}

	vSize := tx.VSize()
	if vSize != 136 {
		t.Errorf("incorrect number of vBytes\nWanted 136\nGot    %d", vSize)
		return
	}
}
