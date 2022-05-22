package der

import (
	"encoding/json"
	"os"
)

type Fixture struct {
	R          string `json:"r"`
	S          string `json:"s"`
	DEREncoded string `json:"derEncoded"`
}

// TODO add more fixtures from
// https://github.com/bitcoin/bitcoin/pull/5713/files
// https://github.com/bitcoinjs/bip66/blob/master/test/fixtures.json

func readFixtures() ([]*Fixture, error) {
	var fixtures []*Fixture

	fixtureData, err := os.ReadFile("fixtures.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(fixtureData, &fixtures); err != nil {
		return nil, err
	}

	return fixtures, nil
}
