package tx

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/bitcoinlib/script"
	"github.com/kklash/bitcoinlib/varint"
)

// Notes
// https://raghavsood.com/blog/2018/06/10/bitcoin-signature-types-sighash
// https://medium.com/@bitaps.com/exploring-bitcoin-signature-hash-types-15427766f0a9
// https://leftasexercise.com/2018/04/06/signing-and-verifying-bitcoin-transactions/
// https://github.com/karask/python-bitcoin-utils/blob/master/bitcoinutils/transactions.py
// ~/github/ExodusMovement/desktop/src/node_modules/bitcoinjs-lib/src/transaction.js

// Clone creates a duplicate of tx.
func (tx *Tx) Clone() *Tx {
	clone := new(Tx)
	clone.Version = tx.Version

	clone.Inputs = make([]*Input, len(tx.Inputs))
	for i, vin := range tx.Inputs {
		if vin != nil {
			clone.Inputs[i] = vin.Clone()
		}
	}

	clone.Outputs = make([]*Output, len(tx.Outputs))
	for i, vout := range tx.Outputs {
		if vout != nil {
			clone.Outputs[i] = vout.Clone()
		}
	}

	if tx.Witnesses != nil {
		clone.Witnesses = make([]Witness, len(tx.Witnesses))
		for i, witness := range tx.Witnesses {
			if witness != nil {
				clone.Witnesses[i] = witness.Clone()
			}
		}
	}

	clone.Locktime = tx.Locktime
	return clone
}

func (tx *Tx) SignatureHashForInput(nInput int, prevOutScript []byte, sigHashType uint32) (hashed [32]byte, err error) {
	tx = tx.Clone()

	sigHashNone := sigHashType&0x1f == constants.SigHashNone
	sigHashSingle := sigHashType&0x1f == constants.SigHashSingle
	sigHashAnyoneCanPay := sigHashType&constants.SigHashAnyoneCanPay > 0

	// https://github.com/bitcoin/bitcoin/blob/master/src/test/sighash_tests.cpp#L29
	// https://github.com/bitcoin/bitcoin/blob/master/src/test/sighash_tests.cpp#L60
	if nInput >= len(tx.Inputs) || (sigHashSingle && nInput >= len(tx.Outputs)) {
		hashed[31] = 1
		return
	}

	tx.Inputs[nInput].Script, err = script.StripOpCode(prevOutScript, constants.OP_CODESEPARATOR)
	if err != nil {
		return
	}

	if sigHashAnyoneCanPay {
		// ignore all other inputs
		tx.Inputs = tx.Inputs[nInput : nInput+1]
	} else {
		for i, vin := range tx.Inputs {
			if i != nInput {
				vin.Script = []byte{}
				if sigHashNone || sigHashSingle {
					vin.Sequence = 0
				}
			}
		}
	}

	if sigHashNone {
		tx.Outputs = tx.Outputs[0:0]
	} else if sigHashSingle {
		// Ignore outputs after the signing input index
		tx.Outputs = tx.Outputs[:nInput+1]

		// Blank outputs before the signing input index
		for i := 0; i < nInput; i++ {
			tx.Outputs[i].Value = 0xffffffffffffffff
			tx.Outputs[i].Script = []byte{}
		}
	}

	hasher := bhash.NewMultiHasher(sha256.New(), sha256.New())
	if _, err = tx.WriteToNoWitness(hasher); err != nil {
		return
	}

	if err = binary.Write(hasher, binary.LittleEndian, sigHashType); err != nil {
		return
	}

	// Note: Signature Hashes are reversed when displayed as hex
	copy(hashed[:], hasher.Sum(nil))

	return
}

func (tx *Tx) hashPrevouts() (hashed [32]byte, err error) {
	hasher := bhash.NewMultiHasher(sha256.New(), sha256.New())

	for _, vin := range tx.Inputs {
		if _, err = vin.PrevOut.WriteTo(hasher); err != nil {
			return
		}
	}

	copy(hashed[:], hasher.Sum(nil))
	return
}

func (tx *Tx) hashSequence() (hashed [32]byte, err error) {
	hasher := bhash.NewMultiHasher(sha256.New(), sha256.New())

	for _, vin := range tx.Inputs {
		if err = binary.Write(hasher, binary.LittleEndian, vin.Sequence); err != nil {
			return
		}
	}

	copy(hashed[:], hasher.Sum(nil))
	return
}

func (tx *Tx) hashAllOutputs() (hashed [32]byte, err error) {
	hasher := bhash.NewMultiHasher(sha256.New(), sha256.New())

	for _, vout := range tx.Outputs {
		if _, err = vout.WriteTo(hasher); err != nil {
			return
		}
	}

	copy(hashed[:], hasher.Sum(nil))
	return
}

func (tx *Tx) hashSingleOutput(n int) (hashed [32]byte, err error) {
	hasher := bhash.NewMultiHasher(sha256.New(), sha256.New())
	if _, err = tx.Outputs[n].WriteTo(hasher); err != nil {
		return
	}

	copy(hashed[:], hasher.Sum(nil))
	return
}

/* https://en.bitcoin.it/wiki/BIP_0143
	 Double SHA256 of the serialization of:
    1. nVersion of the transaction (4-byte little endian)
    2. hashPrevouts (32-byte hash)
    3. hashSequence (32-byte hash)
    4. outpoint (32-byte hash + 4-byte little endian)
    5. scriptCode of the input (serialized as scripts inside CTxOuts)
    6. value of the output spent by this input (8-byte little endian)
    7. nSequence of the input (4-byte little endian)
    8. hashOutputs (32-byte hash)
    9. nLocktime of the transaction (4-byte little endian)
   10. sighash type of the signature (4-byte little endian)
*/

func (tx *Tx) SignatureHashForWitnessInput(nInput int, prevOutScript []byte, sigHashType uint32, inputValue uint64) (hashed [32]byte, err error) {
	var (
		hashPrevouts [32]byte
		hashSequence [32]byte
		hashOutputs  [32]byte

		sigHashNone         = sigHashType&0x1f == constants.SigHashNone
		sigHashSingle       = sigHashType&0x1f == constants.SigHashSingle
		sigHashAnyoneCanPay = sigHashType&constants.SigHashAnyoneCanPay > 0
	)

	if !sigHashAnyoneCanPay {
		hashPrevouts, err = tx.hashPrevouts()
		if err != nil {
			return
		}
	}

	if !sigHashAnyoneCanPay && !sigHashSingle && !sigHashNone {
		hashSequence, err = tx.hashSequence()
		if err != nil {
			return
		}
	}

	if !sigHashSingle && !sigHashNone {
		hashOutputs, err = tx.hashAllOutputs()
	} else if sigHashSingle && nInput < len(tx.Outputs) {
		hashOutputs, err = tx.hashSingleOutput(nInput)
	}
	if err != nil {
		return
	}

	hasher := bhash.NewMultiHasher(sha256.New(), sha256.New())

	if err = binary.Write(hasher, binary.LittleEndian, tx.Version); err != nil {
		return
	}

	if _, err = hasher.Write(hashPrevouts[:]); err != nil {
		return
	}

	if _, err = hasher.Write(hashSequence[:]); err != nil {
		return
	}

	if _, err = tx.Inputs[nInput].PrevOut.WriteTo(hasher); err != nil {
		return
	}

	scriptLen := varint.VarInt(len(prevOutScript))
	if _, err = scriptLen.WriteTo(hasher); err != nil {
		return
	}

	if _, err = hasher.Write(prevOutScript); err != nil {
		return
	}

	if err = binary.Write(hasher, binary.LittleEndian, inputValue); err != nil {
		return
	}

	if err = binary.Write(hasher, binary.LittleEndian, tx.Inputs[nInput].Sequence); err != nil {
		return
	}

	if _, err = hasher.Write(hashOutputs[:]); err != nil {
		return
	}

	if err = binary.Write(hasher, binary.LittleEndian, tx.Locktime); err != nil {
		return
	}

	if err = binary.Write(hasher, binary.LittleEndian, sigHashType); err != nil {
		return
	}

	copy(hashed[:], hasher.Sum(nil))

	return
}
