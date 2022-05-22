package rpc

import (
	"reflect"
	"testing"
)

func TestConnection_RequestSetResult(t *testing.T) {
	t.Skip()

	username, password, err := FindCookie()
	if err != nil {
		t.Errorf("failed to find cookie: %s", err)
		return
	}

	conn, err := NewConnection("http://127.0.0.1:8332", username, password)
	if err != nil {
		t.Errorf("failed to create connection: %s", err)
		return
	}

	var blockCount uint32
	if err := conn.RequestSetResult(&blockCount, "getblockcount"); err != nil {
		t.Errorf("failed to get block count: %s", err)
		return
	} else if blockCount < 100 {
		t.Errorf("Incorrect block count: %d", blockCount)
		return
	}

	type Input struct {
		Txid     string
		Vout     int
		Sequence uint32
	}

	type Output struct {
		Value        float64
		ScriptPubKey struct{ Hex string }
	}

	type Transaction struct {
		Txid     string
		Version  int32
		Size     int
		Locktime uint32
		Inputs   []Input  `json:"vin"`
		Outputs  []Output `json:"vout"`
	}

	var transaction Transaction

	if err := conn.RequestSetResult(&transaction, "getrawtransaction", "3b7f895f7dc9ab91e691755766effa703d0ee9ecb1bae93b2f75540bf0d1e8be", true); err != nil {
		t.Errorf("failed to get transaction: %s", err)
		return
	}

	expectedTransaction := Transaction{
		Txid:     "3b7f895f7dc9ab91e691755766effa703d0ee9ecb1bae93b2f75540bf0d1e8be",
		Version:  2,
		Size:     361,
		Locktime: 684802,
		Inputs: []Input{
			{
				Txid:     "ee18144030174d513425280635e39226958b9ce60e6b974ee0b26bf11c44c92a",
				Vout:     5,
				Sequence: 4294967293,
			},
		},
		Outputs: []Output{
			{
				Value:        0.00114364,
				ScriptPubKey: struct{ Hex string }{"76a914d32a2b0ff8a8eb2471157217d25a11c05b072cfb88ac"},
			},
			{
				Value:        0.00386249,
				ScriptPubKey: struct{ Hex string }{"76a914f80d5ff4796ecfae2df1293031ee6e624b6b379d88ac"},
			},
			{
				Value:        0.00718372,
				ScriptPubKey: struct{ Hex string }{"76a9147862663debf40dbc0a4be0f72626782a0cbdf08088ac"},
			},
			{
				Value:        0.07499346,
				ScriptPubKey: struct{ Hex string }{"76a9141f55642a6719fcaeac910975bc92777738d5c0a088ac"},
			},
			{
				Value:        0.13439691,
				ScriptPubKey: struct{ Hex string }{"76a91490e3704d1a3636b1831661deff10decc25cd235488ac"},
			},
			{
				Value:        17.81134414,
				ScriptPubKey: struct{ Hex string }{"001470da7a4dac44f36361e95c8c5f6419576f27f437"},
			},
		},
	}

	if !reflect.DeepEqual(&transaction, &expectedTransaction) {
		t.Errorf("failed to parse transaction correctly")
		return
	}

	result, err := conn.Request("getblockheader", "000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506")
	if err != nil {
		t.Errorf("Failed to fetch block header: %s", err)
		return
	}

	expectedBlockHeader := map[string]any{
		"bits":              "1b04864c",
		"chainwork":         "0000000000000000000000000000000000000000000000000644cb7f5234089e",
		"confirmations":     620008.0,
		"difficulty":        14484.1623612254,
		"hash":              "000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506",
		"height":            100000.0,
		"mediantime":        1293622620.0,
		"merkleroot":        "f3e94742aca4b5ef85488dc37c06c3282295ffec960994b2c0d5ac2a25a95766",
		"nTx":               4.0,
		"nextblockhash":     "00000000000080b66c911bd5ba14a74260057311eaeb1982802f7010f1a9f090",
		"nonce":             274148111.0,
		"previousblockhash": "000000000002d01c1fccc21636b607dfd930d31d01c3a62104612a1719011250",
		"time":              1293623863.0,
		"version":           1.0,
		"versionHex":        "00000001",
	}

	if !reflect.DeepEqual(result, expectedBlockHeader) {
		t.Errorf("failed to parse block header correctly")
		return
	}
}

// func TestConnection_Request(t *testing.T) {
// 		username, password, err := FindCookie()
// 	if err != nil {
// 		t.Errorf("failed to find cookie: %s", err)
// 		return
// 	}

// 	conn, err := NewConnection("http://127.0.0.1:8332", username, password)
// 	if err != nil {
// 		t.Errorf("failed to create connection: %s", err)
// 		return
// 	}

// }
