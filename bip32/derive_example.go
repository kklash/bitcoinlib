package bip32

import (
	"encoding/hex"
	"fmt"
)

func ExampleDerivePrivateChild() {
	parentKey, _ := hex.DecodeString("f1c7c871a54a804afe328b4c83a1c33b8e5ff48f5087273f04efa83b247d6a2d")
	parentChainCode, _ := hex.DecodeString("637807030d55d01f9a0cb3a7839515d796bd07706386a6eddf06cc29a65a0e29")

	childKey, childChainCode := DerivePrivateChild(parentKey, parentChainCode, 2)
	fmt.Printf("child key: %x\n", childKey)
	fmt.Printf("child chain code: %x\n", childChainCode)

	// Output:
	// child key: bb7d39bdb83ecf58f2fd82b6d918341cbef428661ef01ab97c28a4842125ac23
	// child chain code: 9452b549be8cea3ecb7a84bec10dcfd94afe4d129ebfd3b3cb58eedf394ed271
}

func ExampleDerivePublicChild() {
	parentKey, _ := hex.DecodeString("02d2b36900396c9282fa14628566582f206a5dd0bcc8d5e892611806cafb0301f0")
	parentChainCode, _ := hex.DecodeString("637807030d55d01f9a0cb3a7839515d796bd07706386a6eddf06cc29a65a0e29")

	childKey, childChainCode, err := DerivePublicChild(parentKey, parentChainCode, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("child key: %x\n", childKey)
	fmt.Printf("child chain code: %x\n", childChainCode)

	// Output:
	// child key: 024d902e1a2fc7a8755ab5b694c575fce742c48d9ff192e63df5193e4c7afe1f9c
	// child chain code: 9452b549be8cea3ecb7a84bec10dcfd94afe4d129ebfd3b3cb58eedf394ed271
}
