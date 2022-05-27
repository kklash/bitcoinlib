package bip38

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestEncryptECMultiply(t *testing.T) {
	// TODO why do the official test vectors use incorrect P2PKH addresses?
	t.Skip()

	type Fixture struct {
		password         string
		encryptedKey     string
		privateKey       string
		address          string
		intermediateCode string
		entropy          string
		seedb            string
		compressed       bool
		lot, sequence    uint32
	}

	fixtures := []*Fixture{
		{
			password:         "TestingOneTwoThree",
			encryptedKey:     "6PfQu77ygVyJLZjfvMLyhLMQbYnu5uguoJJ4kMCLqWwPEdfpwANVS76gTX",
			privateKey:       "A43A940577F4E97F5C4D39EB14FF083A98187C64EA7C99EF7CE460833959A519",
			address:          "1PE6TQi6HTVNz5DLwB1LcpMBALubfuN2z2",
			intermediateCode: "passphrasepxFy57B9v8HtUsszJYKReoNDV6VHjUSGt8EVJmux9n1J3Ltf1gRxyDGXqnf9qm",
			entropy:          "a50dba6772cb9383",
			seedb:            "99241d58245c883896f80843d2846672d7312e6195ca1a6c",
			compressed:       false,
		},
		{
			password:         "Satoshi",
			encryptedKey:     "6PfLGnQs6VZnrNpmVKfjotbnQuaJK4KZoPFrAjx1JMJUa1Ft8gnf5WxfKd",
			privateKey:       "C2C8036DF268F498099350718C4A3EF3984D2BE84618C2650F5171DCC5EB660A",
			address:          "1CqzrtZC6mXSAhoxtFwVjz8LtwLJjDYU3V",
			intermediateCode: "passphraseoRDGAXTWzbp72eVbtUDdn1rwpgPUGjNZEc6CGBo8i5EC1FPW8wcnLdq4ThKzAS",
			entropy:          "67010a9573418906",
			seedb:            "49111e301d94eab339ff9f6822ee99d9f49606db3b47a497",
			compressed:       false,
		},
		{
			password:         "MOLON LABE",
			encryptedKey:     "6PgNBNNzDkKdhkT6uJntUXwwzQV8Rr2tZcbkDcuC9DZRsS6AtHts4Ypo1j",
			privateKey:       "44EA95AFBF138356A05EA32110DFD627232D0F2991AD221187BE356F19FA8190",
			address:          "1Jscj8ALrYu2y9TD8NrpvDBugPedmbj4Yh",
			intermediateCode: "passphraseaB8feaLQDENqCgr4gKZpmf4VoaT6qdjJNJiv7fsKvjqavcJxvuR1hy25aTu5sX",
			entropy:          "4fca5a97",
			seedb:            "87a13b07858fa753cd3ab3f1c5eafb5f12579b6c33c9a53f",
			compressed:       false,
			lot:              263183,
			sequence:         1,
		},
		{
			password:         "ΜΟΛΩΝ ΛΑΒΕ",
			encryptedKey:     "6PgGWtx25kUg8QWvwuJAgorN6k9FbE25rv5dMRwu5SKMnfpfVe5mar2ngH",
			privateKey:       "CA2759AA4ADB0F96C414F36ABEB8DB59342985BE9FA50FAAC228C8E7D90E3006",
			address:          "1Lurmih3KruL4xDB5FmHof38yawNtP9oGf",
			intermediateCode: "passphrased3z9rQJHSyBkNBwTRPkUGNVEVrUAcfAXDyRU1V28ie6hNFbqDwbFBvsTK7yWVK",
			entropy:          "c40ea76f",
			seedb:            "03b06a1ea7f9219ae364560d7b985ab1fa27025aaa7e427a",
			compressed:       false,
			lot:              806938,
			sequence:         1,
		},
	}

	for _, fixture := range fixtures {
		entropy, err := hex.DecodeString(fixture.entropy)
		if err != nil {
			t.Errorf("failed to decode fixture entropy %q: %s", fixture.entropy, err)
			continue
		}
		randReader := bytes.NewReader(entropy)

		var intermediateCode string

		if fixture.lot == 0 {
			intermediateCode, err = GenerateIntermediateCode(randReader, fixture.password)
		} else {
			intermediateCode, err = GenerateIntermediateCodeWithLotSequence(
				randReader,
				fixture.password,
				fixture.lot,
				fixture.sequence,
			)
		}

		if err != nil {
			t.Errorf("failed to generate intermediate code %s: %s", fixture.intermediateCode, err)
			continue
		}
		if intermediateCode != fixture.intermediateCode {
			t.Errorf("generated incorrect intermediate code\nWanted %s\nGot    %s", fixture.intermediateCode, intermediateCode)
			continue
		}

		privateKey, compressed, err := Decrypt(fixture.encryptedKey, fixture.password)
		if err != nil {
			t.Errorf("failed to decrypt ecmult private key %s: %s", fixture.encryptedKey, err)
			continue
		}

		if compressed != fixture.compressed {
			t.Errorf("expected to receive compressed=%v for key %s", fixture.compressed, fixture.encryptedKey)
			continue
		}

		expectedPrivateKey, err := hex.DecodeString(fixture.privateKey)
		if err != nil {
			t.Errorf("failed to decode fixture private key hex %q: %s", fixture.privateKey, err)
			continue
		}

		if !bytes.Equal(privateKey, expectedPrivateKey) {
			t.Errorf("decrypted private key does not match\nWanted %X\nGot    %X", expectedPrivateKey, privateKey)
			continue
		}

		seedb, err := hex.DecodeString(fixture.seedb)
		if err != nil {
			t.Errorf("failed to decode fixture seedb %q: %s", fixture.seedb, err)
			continue
		}

		encryptedKey, err := EncryptIntermediateCode(bytes.NewReader(seedb), intermediateCode, fixture.compressed)
		if err != nil {
			t.Errorf("failed to encrypt intermediate code %q: %s", intermediateCode, err)
			continue
		}

		if encryptedKey != fixture.encryptedKey {
			t.Errorf("failed to generate expected encrypted ecmult key\nWanted %s\nGot    %s", fixture.encryptedKey, encryptedKey)
			continue
		}
	}
}
