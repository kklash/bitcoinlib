package unspent

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/kklash/bitcoinlib/blocks"
	"github.com/kklash/bitcoinlib/tx"
)

func mustHex(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return decoded
}

func mustHash(s string) (hash [32]byte) {
	decoded := mustHex(s)
	copy(hash[:], decoded)
	return
}

func TestUpdateFromBlock(t *testing.T) {
	// trimmed down version of block 684836, with only the relevant transactions.
	blockRaw := mustHex("04e0ff27193abf9d8ba5883e0de9f65040e34a1c8bb13fb4f404040000000000000000008c6d1a33788da790d1935d76a0fd17aa1b87f9fd797249b2813d2ac733c54e30a43bac60e93c0b171bfdf69102020000000001012ac9441cf16bb2e04e976b0ee69c8b952692e33506282534514d1730401418ee0500000000fdffffff06bcbe0100000000001976a914d32a2b0ff8a8eb2471157217d25a11c05b072cfb88acc9e40500000000001976a914f80d5ff4796ecfae2df1293031ee6e624b6b379d88ac24f60a00000000001976a9147862663debf40dbc0a4be0f72626782a0cbdf08088ac526e7200000000001976a9141f55642a6719fcaeac910975bc92777738d5c0a088accb12cd00000000001976a91490e3704d1a3636b1831661deff10decc25cd235488ac4ef4296a0000000016001470da7a4dac44f36361e95c8c5f6419576f27f4370247304402201098f331b780282841760a1a2a7fc9c2683c0411285e387b625a48c70e063dc402202263cd3c5dde4a013d1ef8bf37b865625fd2ae25859fa91fc2f9ad9041a3b3b3012103f392fd19651faa764b75583f5fe6bfcd633f1d7be6256000f79e7240789ff74002730a0002000000000101bee8d1f00b54752f3be9bab1ece90e3d70faef66577591e691abc97d5f897f3b0500000000fdffffff0220770e00000000001976a914d32a2b0ff8a8eb2471157217d25a11c05b072cfb88aca6ec1a6a00000000160014bdb7898232932076327ba62aa2f55ee720f358d802483045022100ff276ecf1751657ab44ac38a6af8ae922f0b95cffb2051618acc12cc518ecef90220227ef68d71c05dbf66d3fe6652939846e70a4252c4a3990ae2b57cfaf5c24605012103615aad518b122eb8238d692e0f9a96178245568aefe086181a6b7724f3564db323730a00")
	block, _ := blocks.FromReader(bytes.NewReader(blockRaw))

	scriptPubKeys := [][]byte{
		mustHex("001481916950b977370407e58cb3970ce1292093e6e6"),       // random non-existent address
		mustHex("00148dc8e4f2450a1af2202383d71a259de6ea14d4ba"),       // bc1q3hywfuj9pgd0ygprs0t35fvaum4pf496etuasm
		mustHex("76a914d32a2b0ff8a8eb2471157217d25a11c05b072cfb88ac"), // 1LFY6NkSx2XdtLY8w1a8pv8ruMnprGM8eE
	}

	utxos := []*Output{
		{
			// Nonexistent utxo, expect it to be left untouched
			Outpoint: &tx.PrevOut{
				Hash:  mustHash("82a6e289b344ae3cd67350427df5133c75c0d8fd3c6f5d342aaeb8817b807ac5"),
				Index: 2,
			},
			TxOut: &tx.Output{
				Value:  100000,
				Script: scriptPubKeys[0],
			},
		},
		{
			// Pre-existing utxo which is spent, expect it to be deleted from the output set
			Outpoint: &tx.PrevOut{
				Hash:  mustHash("2ac9441cf16bb2e04e976b0ee69c8b952692e33506282534514d1730401418ee"),
				Index: 5,
			},
			TxOut: &tx.Output{
				Value:  1803477436,
				Script: scriptPubKeys[1],
			},
		},
	}

	unspentOutputs := new(OutputSet)
	for _, utxo := range utxos {
		unspentOutputs.AddOutput(utxo)
	}

	if err := unspentOutputs.UpdateFromBlock(block, scriptPubKeys); err != nil {
		t.Errorf("failed to update unspent output set: %s", err)
		return
	}

	expectedOutputs := []*Output{
		utxos[0], // Non-existent utxo should be untouched
		{
			// New UTXO sent to 1LFY6NkSx2XdtLY8w1a8pv8ruMnprGM8eE
			Outpoint: &tx.PrevOut{
				Hash:  mustHash("3691e0552f263c21adbe60accb38a493d037f47bc122085840b4b1b30f32e245"),
				Index: 0,
			},
			TxOut: &tx.Output{
				Value:  948000,
				Script: scriptPubKeys[2],
			},
		},
		{
			// New UTXO sent to 1LFY6NkSx2XdtLY8w1a8pv8ruMnprGM8eE
			Outpoint: &tx.PrevOut{
				Hash:  mustHash("bee8d1f00b54752f3be9bab1ece90e3d70faef66577591e691abc97d5f897f3b"),
				Index: 0,
			},
			TxOut: &tx.Output{
				Value:  114364,
				Script: scriptPubKeys[2],
			},
		},
	}

	if unspentOutputs.Size() != len(expectedOutputs) {
		t.Errorf("unexpected number of resulting unspent outputs\nWanted %d\nGot    %d", len(expectedOutputs), unspentOutputs.Size())
		return
	}

	for _, expectedUtxo := range expectedOutputs {
		actualUtxo := unspentOutputs.GetByHash(expectedUtxo.Outpoint.Hash, expectedUtxo.Outpoint.Index)
		if actualUtxo == nil {
			t.Errorf("expected UTXO missing:\n hash: %x\n index: %d", expectedUtxo.Outpoint.Hash, expectedUtxo.Outpoint.Index)
			continue
		}

		if !reflect.DeepEqual(actualUtxo, expectedUtxo) {
			t.Errorf("utxo does not match expected:\n hash: %x\n index: %d", expectedUtxo.Outpoint.Hash, expectedUtxo.Outpoint.Index)
			continue
		}
	}
}
