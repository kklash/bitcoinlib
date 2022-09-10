package script

import (
	"bytes"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/bitcoinlib/taproot"
)

var (
	taprootLeafHasher   = bhash.NewTaggedHasher("TapLeaf")
	taprootBranchHasher = bhash.NewTaggedHasher("TapBranch")
)

// Hasher is used as a union type for the nodes of a merkelized abstract syntax tree (MAST).
// Each node of the tree must be able to hash itself or its child nodes if applicable.
//
// MastBranch is used for branch nodes which have two children. MastLeaf is used for
// visible script leaf nodes, which are usable as spending conditions for a taproot output.
// MastLeafHash is used for leaves of the MAST which are already hashed, whose script
// preimage is not known.
//
// Instances of these types can be arranged in a binary tree structure like so.
//
//	var mast Hasher = MastBranch{
//		&MastLeaf{
//			Version: 0xc0,
//			Script:  []byte{},
//		},
//		MastBranch{
//			MastLeafHash{1, 2, 3},
//			&MastLeaf{
//				Version: 0xc0,
//				Script:  []byte{},
//			},
//		},
//	}
type Hasher interface {
	Hash() [32]byte
}

// MastLeaf stores a MAST script and version number which can be hashed as part of
// a merkelized abstract syntax tree.
type MastLeaf struct {
	Version byte
	Script  []byte
}

// Hash hashes the MastLeaf version number and push-data prepended script.
func (ms *MastLeaf) Hash() (hashed [32]byte) {
	pushScript := PushData(ms.Script)
	preimage := make([]byte, 1+len(pushScript))
	preimage[0] = ms.Version
	copy(preimage[1:], pushScript)
	copy(hashed[:], taprootLeafHasher(preimage))
	return
}

// MastLeafHash represents a node on a merkelized abstract syntax tree that
// has already been hashed. Using nodes of this kind allow callers to obfuscate
// unused locking script conditions.
type MastLeafHash [32]byte

// Hash returns the MastLeafHash as a [32]byte array.
func (hashed MastLeafHash) Hash() [32]byte {
	return hashed
}

// MastBranch represents a branch node in a merkelized abstract syntax tree with two child nodes,
// which must satisfy type Hasher.
type MastBranch [2]Hasher

// Hash hashes both branches of the MastBranch and returns their hash. Panics if either
// the left or the right branch is nil.
func (mn MastBranch) Hash() (hashed [32]byte) {
	if mn[0] == nil || mn[1] == nil {
		panic("MastBranch is missing a node - MAST trees should always have two child elements")
	}
	leftH := mn[0].Hash()
	rightH := mn[1].Hash()

	// Sorting ensures the tree is deterministic regardless of how branches are arranged WRT left vs right
	if bytes.Compare(leftH[:], rightH[:]) == 1 {
		leftH, rightH = rightH, leftH
	}
	copy(hashed[:], taprootBranchHasher(append(leftH[:], rightH[:]...)))
	return
}

// MakeP2TR creates a taproot P2TR output from the given internal public key
// and optional MAST script tree. If scriptTree is nil, the unspendable empty
// script is used in the taproot script commitment. This means the only way
// to redeem the output would be to provide a signature by internalPublicKey.
func MakeP2TR(internalPublicKey []byte, scriptTree Hasher) ([]byte, error) {
	h := []byte{}
	if scriptTree != nil {
		scriptTreeHash := scriptTree.Hash()
		h = scriptTreeHash[:]
	}
	outputPublicKey, _, err := taproot.TweakPublicKey(internalPublicKey, h)
	if err != nil {
		return nil, err
	}

	scriptPubKey := append([]byte{constants.OP_TRUE}, PushData(outputPublicKey)...)
	return scriptPubKey, nil
}
