package bip32

import (
	"encoding/hex"
	"testing"

	"github.com/kklash/bitcoinlib/constants"
)

type NodeFixture struct {
	xpub     string
	xpriv    string
	children map[uint32]*NodeFixture
}

func (nf *NodeFixture) Test(t *testing.T) {
	privateKey, chainCode, parentFingerprint, depth, index, privateVersion, err := Deserialize(nf.xpriv)
	if err != nil {
		t.Errorf("Failed to deserialize xpriv key: %s\nError: %s", nf.xpriv, err)
		return
	}

	publicKey, _, _, _, _, publicVersion, err := Deserialize(nf.xpub)
	if err != nil {
		t.Errorf("Failed to deserialize xpub key: %s\nError: %s", nf.xpub, err)
		return
	}

	xpriv := SerializePrivate(privateKey, chainCode, parentFingerprint, depth, index, privateVersion)
	if xpriv != nf.xpriv {
		t.Errorf("Generated xpriv did not match fixture\n wanted %s\n got %s", nf.xpriv, xpriv)
		return
	}

	xpub := SerializePublic(publicKey, chainCode, parentFingerprint, depth, index, publicVersion)
	if xpub != nf.xpub {
		t.Errorf("Generated xpub did not match fixture\n wanted %s\n got %s", nf.xpub, xpub)
		return
	}

	fingerprint, err := KeyFingerprint(publicKey)
	if err != nil {
		t.Errorf("Failed to generate key fingerprint for %x\nError: %s", publicKey, err)
		return
	}

	if nf.children == nil {
		return
	}

	for index, cf := range nf.children {
		childPrivateKey, childChainCode := DerivePrivateChild(privateKey, chainCode, index)
		var childPublicKey []byte
		if index < constants.Bip32Hardened {
			if childPublicKey, _, err = DerivePublicChild(publicKey, chainCode, index); err != nil {
				t.Errorf("Failed to derive child public key: %s\nError: %s", cf.xpub, err)
				return
			}
		} else {
			childPublicKey = NeuterCompressed(childPrivateKey)
		}

		childXpub := SerializePublic(childPublicKey, childChainCode, fingerprint, depth+1, index, publicVersion)
		if childXpub != cf.xpub {
			t.Errorf("Failed to derive correct child public key\n wanted %s\n got %s", cf.xpub, childXpub)
			return
		}

		childXpriv := SerializePrivate(childPrivateKey, childChainCode, fingerprint, depth+1, index, privateVersion)
		if childXpriv != cf.xpriv {
			t.Errorf("Failed to derive correct child private key\n wanted %s\n got %s", cf.xpriv, childXpriv)
			return
		}

		cf.Test(t)
	}
}

type MasterFixture struct {
	seedHex string
	node    *NodeFixture
}

func (mf *MasterFixture) Test(t *testing.T) {
	seed, _ := hex.DecodeString(mf.seedHex)
	masterKey, chainCode, err := GenerateMasterKey(seed)
	if err != nil {
		t.Errorf("Failed to generate master key: %s\nError: %s", mf.node.xpriv, err)
		return
	}

	xpriv := SerializePrivate(masterKey, chainCode, nil, 0, 0, constants.BitcoinNetwork.ExtendedPrivate)
	if xpriv != mf.node.xpriv {
		t.Errorf("Failed to generate correct master private key\n wanted %s\ngot %s", mf.node.xpriv, xpriv)
		return
	}

	masterPublicKey := NeuterCompressed(masterKey)

	xpub := SerializePublic(masterPublicKey, chainCode, nil, 0, 0, constants.BitcoinNetwork.ExtendedPublic)
	if xpub != mf.node.xpub {
		t.Errorf("Failed to generate correct master public key\n wanted %s\ngot %s", mf.node.xpub, xpub)
		return
	}

	mf.node.Test(t)
}

func TestBip32Vectors(t *testing.T) {
	// Test vectors from https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
	var fixtures = []MasterFixture{
		{
			seedHex: "000102030405060708090a0b0c0d0e0f",
			node: &NodeFixture{
				xpub:  "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8",
				xpriv: "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi",
				children: map[uint32]*NodeFixture{
					0x80000000: {
						xpub:  "xpub68Gmy5EdvgibQVfPdqkBBCHxA5htiqg55crXYuXoQRKfDBFA1WEjWgP6LHhwBZeNK1VTsfTFUHCdrfp1bgwQ9xv5ski8PX9rL2dZXvgGDnw",
						xpriv: "xprv9uHRZZhk6KAJC1avXpDAp4MDc3sQKNxDiPvvkX8Br5ngLNv1TxvUxt4cV1rGL5hj6KCesnDYUhd7oWgT11eZG7XnxHrnYeSvkzY7d2bhkJ7",
						children: map[uint32]*NodeFixture{
							1: {
								xpub:  "xpub6ASuArnXKPbfEwhqN6e3mwBcDTgzisQN1wXN9BJcM47sSikHjJf3UFHKkNAWbWMiGj7Wf5uMash7SyYq527Hqck2AxYysAA7xmALppuCkwQ",
								xpriv: "xprv9wTYmMFdV23N2TdNG573QoEsfRrWKQgWeibmLntzniatZvR9BmLnvSxqu53Kw1UmYPxLgboyZQaXwTCg8MSY3H2EU4pWcQDnRnrVA1xe8fs",
								children: map[uint32]*NodeFixture{
									0x80000002: {
										xpub:  "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
										xpriv: "xprv9z4pot5VBttmtdRTWfWQmoH1taj2axGVzFqSb8C9xaxKymcFzXBDptWmT7FwuEzG3ryjH4ktypQSAewRiNMjANTtpgP4mLTj34bhnZX7UiM",
										children: map[uint32]*NodeFixture{
											2: {
												xpub:  "xpub6FHa3pjLCk84BayeJxFW2SP4XRrFd1JYnxeLeU8EqN3vDfZmbqBqaGJAyiLjTAwm6ZLRQUMv1ZACTj37sR62cfN7fe5JnJ7dh8zL4fiyLHV",
												xpriv: "xprvA2JDeKCSNNZky6uBCviVfJSKyQ1mDYahRjijr5idH2WwLsEd4Hsb2Tyh8RfQMuPh7f7RtyzTtdrbdqqsunu5Mm3wDvUAKRHSC34sJ7in334",
												children: map[uint32]*NodeFixture{
													1000000000: {
														xpub:  "xpub6H1LXWLaKsWFhvm6RVpEL9P4KfRZSW7abD2ttkWP3SSQvnyA8FSVqNTEcYFgJS2UaFcxupHiYkro49S8yGasTvXEYBVPamhGW6cFJodrTHy",
														xpriv: "xprvA41z7zogVVwxVSgdKUHDy1SKmdb533PjDz7J6N6mV6uS3ze1ai8FHa8kmHScGpWmj4WggLyQjgPie1rFSruoUihUZREPSL39UNdE3BBDu76",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		{
			seedHex: "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542",
			node: &NodeFixture{
				xpub:  "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB",
				xpriv: "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U",
				children: map[uint32]*NodeFixture{
					0: {
						xpub:  "xpub69H7F5d8KSRgmmdJg2KhpAK8SR3DjMwAdkxj3ZuxV27CprR9LgpeyGmXUbC6wb7ERfvrnKZjXoUmmDznezpbZb7ap6r1D3tgFxHmwMkQTPH",
						xpriv: "xprv9vHkqa6EV4sPZHYqZznhT2NPtPCjKuDKGY38FBWLvgaDx45zo9WQRUT3dKYnjwih2yJD9mkrocEZXo1ex8G81dwSM1fwqWpWkeS3v86pgKt",
						children: map[uint32]*NodeFixture{
							(0x80000000 + 2147483647): {
								xpub:  "xpub6ASAVgeehLbnwdqV6UKMHVzgqAG8Gr6riv3Fxxpj8ksbH9ebxaEyBLZ85ySDhKiLDBrQSARLq1uNRts8RuJiHjaDMBU4Zn9h8LZNnBC5y4a",
								xpriv: "xprv9wSp6B7kry3Vj9m1zSnLvN3xH8RdsPP1Mh7fAaR7aRLcQMKTR2vidYEeEg2mUCTAwCd6vnxVrcjfy2kRgVsFawNzmjuHc2YmYRmagcEPdU9",
								children: map[uint32]*NodeFixture{
									1: {
										xpub:  "xpub6DF8uhdarytz3FWdA8TvFSvvAh8dP3283MY7p2V4SeE2wyWmG5mg5EwVvmdMVCQcoNJxGoWaU9DCWh89LojfZ537wTfunKau47EL2dhHKon",
										xpriv: "xprv9zFnWC6h2cLgpmSA46vutJzBcfJ8yaJGg8cX1e5StJh45BBciYTRXSd25UEPVuesF9yog62tGAQtHjXajPPdbRCHuWS6T8XA2ECKADdw4Ef",
										children: map[uint32]*NodeFixture{
											(0x80000000 + 2147483646): {
												xpub:  "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
												xpriv: "xprvA1RpRA33e1JQ7ifknakTFpgNXPmW2YvmhqLQYMmrj4xJXXWYpDPS3xz7iAxn8L39njGVyuoseXzU6rcxFLJ8HFsTjSyQbLYnMpCqE2VbFWc",
												children: map[uint32]*NodeFixture{
													2: {
														xpub:  "xpub6FnCn6nSzZAw5Tw7cgR9bi15UV96gLZhjDstkXXxvCLsUXBGXPdSnLFbdpq8p9HmGsApME5hQTZ3emM2rnY5agb9rXpVGyy3bdW6EEgAtqt",
														xpriv: "xprvA2nrNbFZABcdryreWet9Ea4LvTJcGsqrMzxHx98MMrotbir7yrKCEXw7nadnHM8Dq38EGfSh6dqA9QWTyefMLEcBYJUuekgW4BYPJcr9E7j",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		{
			seedHex: "4b381541583be4423346c643850da4b320e46a87ae3d2a4e6da11eba819cd4acba45d239319ac14f863b8d5ab5a0d0c64d2e8a1e7d1457df2e5a3c51c73235be",
			node: &NodeFixture{
				xpub:  "xpub661MyMwAqRbcEZVB4dScxMAdx6d4nFc9nvyvH3v4gJL378CSRZiYmhRoP7mBy6gSPSCYk6SzXPTf3ND1cZAceL7SfJ1Z3GC8vBgp2epUt13",
				xpriv: "xprv9s21ZrQH143K25QhxbucbDDuQ4naNntJRi4KUfWT7xo4EKsHt2QJDu7KXp1A3u7Bi1j8ph3EGsZ9Xvz9dGuVrtHHs7pXeTzjuxBrCmmhgC6",
				children: map[uint32]*NodeFixture{
					0x80000000: {
						xpub:  "xpub68NZiKmJWnxxS6aaHmn81bvJeTESw724CRDs6HbuccFQN9Ku14VQrADWgqbhhTHBaohPX4CjNLf9fq9MYo6oDaPPLPxSb7gwQN3ih19Zm4Y",
						xpriv: "xprv9uPDJpEQgRQfDcW7BkF7eTya6RPxXeJCqCJGHuCJ4GiRVLzkTXBAJMu2qaMWPrS7AANYqdq6vcBcBUdJCVVFceUvJFjaPdGZ2y9WACViL4L",
					},
				},
			},
		},
	}

	for _, fixture := range fixtures {
		fixture.Test(t)
	}
}
