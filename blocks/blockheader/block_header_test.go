package blockheader

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func decodeHex256(hexStr string) (h256 [32]byte, err error) {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return
	}

	copy(h256[:], decoded)
	return
}

type Fixture struct {
	Hex                   string `json:"raw"`
	Hash                  string `json:"hash"`
	Version               int32  `json:"version"`
	PreviousHeaderHashHex string `json:"previousblockhash"`
	MerkleRootHex         string `json:"merkleroot"`
	NBits                 uint32 `json:"bits"`
	TargetNBits           string `json:"target_nbits"`
	Nonce                 uint32 `json:"nonce"`
	Time                  uint32 `json:"time"`
}

func (f *Fixture) GetHeaderStruct() (*BlockHeader, error) {
	header := new(BlockHeader)
	header.Time = f.Time
	header.Nonce = f.Nonce
	header.NBits = f.NBits

	var err error
	if header.MerkleRootHash, err = decodeHex256(f.MerkleRootHex); err != nil {
		return nil, err
	}
	if header.PreviousHeaderHash, err = decodeHex256(f.PreviousHeaderHashHex); err != nil {
		return nil, err
	}

	return header, nil
}

func getFixtures() ([]*Fixture, error) {
	fixtureBytes, err := ioutil.ReadFile("fixtures.json")
	if err != nil {
		return nil, err
	}

	var fixtures []*Fixture
	if err := json.Unmarshal(fixtureBytes, &fixtures); err != nil {
		return nil, err
	}

	return fixtures, nil
}

func testBlockHeader() error {
	fixtures, err := getFixtures()
	if err != nil {
		return err
	}

	for _, fixture := range fixtures {
		rawBytes, _ := hex.DecodeString(fixture.Hex)
		header, err := FromReader(bytes.NewBuffer(rawBytes))
		if err != nil {
			return err
		}

		if header.Version != fixture.Version {
			return fmt.Errorf("Version does not match\nWanted %d\nGot    %d", fixture.Version, header.Version)
		} else if actual := hex.EncodeToString(header.PreviousHeaderHash[:]); actual != fixture.PreviousHeaderHashHex {
			return fmt.Errorf("Previous Header Hash does not match\nWanted %s\nGot    %s", fixture.PreviousHeaderHashHex, actual)
		} else if header.Nonce != fixture.Nonce {
			return fmt.Errorf("Nonce does not match\nWanted %d\nGot    %d", fixture.Nonce, header.Nonce)
		} else if header.Time != fixture.Time {
			return fmt.Errorf("Time does not match\nWanted %d\nGot    %d", fixture.Time, header.Time)
		} else if header.NBits != fixture.NBits {
			return fmt.Errorf("NBits does not match\nWanted %d\nGot    %d", fixture.NBits, header.NBits)
		}

		headerHash, err := header.Hash()
		if err != nil {
			return err
		} else if headerHashHex := hex.EncodeToString(headerHash[:]); headerHashHex != fixture.Hash {
			return fmt.Errorf("Header hash does not match\nWanted %s\nGot    %s", fixture.Hash, headerHashHex)
		}

		if nBitsHex := fmt.Sprintf("%x", header.TargetNBits()); nBitsHex != fixture.TargetNBits {
			return fmt.Errorf("Target NBits does not match\nWanted %s\nGot    %s", fixture.TargetNBits, nBitsHex)
		}

		buf := new(bytes.Buffer)

		n, err := header.WriteTo(buf)
		if err != nil {
			return fmt.Errorf("Failed to encode block header: %s", err)
		} else if n != BlockHeaderSize {
			return fmt.Errorf("unexpected block header bytes written count\nWanted %d\nGot    %d", BlockHeaderSize, n)
		}

		if blockHeaderHex := hex.EncodeToString(buf.Bytes()); blockHeaderHex != fixture.Hex {
			return fmt.Errorf("unexpected block header encoding\nWanted %s\nGot    %s", fixture.Hex, blockHeaderHex)
		}
	}

	return nil
}

func TestBlockHeader(t *testing.T) {
	if err := testBlockHeader(); err != nil {
		t.Error(err)
	}
}
