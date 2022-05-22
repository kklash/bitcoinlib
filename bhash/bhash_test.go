package bhash

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func hex2bytes(h string) []byte {
	buf, _ := hex.DecodeString(h)
	return buf
}

func runHashTest(t *testing.T, hashFn func([]byte) []byte, inputData, outputData [][]byte) {
	for i := 0; i < len(inputData); i++ {
		hashed := hashFn(inputData[i])
		if !bytes.Equal(hashed, outputData[i]) {
			t.Errorf("hash for data fixture does not match - %x", inputData[i])
		}
	}
}

func TestSha256(t *testing.T) {
	inputData := [][]byte{
		hex2bytes("01020304"),
		hex2bytes("deadbeef"),
	}

	outputData := [][]byte{
		hex2bytes("9f64a747e1b97f131fabb6b447296c9b6f0201e79fb3c5356e6c77e89b6a806a"),
		hex2bytes("5f78c33274e43fa9de5659265c1d917e25c03722dcb0b8d27db8d5feaa813953"),
	}

	hashFn := func(d []byte) []byte {
		hash := Sha256(d)
		return hash[:]
	}
	runHashTest(t, hashFn, inputData, outputData)
}

func TestHash160(t *testing.T) {
	inputData := [][]byte{
		[]byte("my foobar string"),
		hex2bytes("00000000000000000021111111111111111112"),
	}

	outputData := [][]byte{
		hex2bytes("734021b5fffef703cb6f7e01a782788a0a2ec0ce"),
		hex2bytes("90cc2eb8a67ab78c15c18a99c1d0ef6b29a990f9"),
	}

	hashFn := func(d []byte) []byte {
		hash := Hash160(d)
		return hash[:]
	}
	runHashTest(t, hashFn, inputData, outputData)
}

func TestDoubleSha256(t *testing.T) {
	inputData := [][]byte{
		[]byte("        oi what u doin"),
		[]byte{9, 8, 7, 6, 5, 4, 3, 2, 1},
		hex2bytes("beefbaaa"),
	}

	outputData := [][]byte{
		hex2bytes("feb97b89aeec3c496751c832edf87204bf22c78b22c7ab71f5f5c118c672aa7a"),
		hex2bytes("df8a6fab8aafae1984ebb958a398164c82b0a119fcbb1fff80c3e577a71b972b"),
		hex2bytes("4490a7e69e7a7aab68c1a8298211989a3dfa699039949e85035b743736911505"),
	}

	hashFn := func(d []byte) []byte {
		hash := DoubleSha256(d)
		return hash[:]
	}
	runHashTest(t, hashFn, inputData, outputData)
}
