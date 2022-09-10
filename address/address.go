// Package address provides APIs to encode and decode addresses for bitcoin-like coins.
package address

import (
	"fmt"

	"github.com/kklash/bitcoinlib/base58check"
	"github.com/kklash/bitcoinlib/bech32"
	"github.com/kklash/bitcoinlib/constants"
	"github.com/kklash/bitcoinlib/script"
)

// Make can generate an address of any given AddressFormat. For example, if you want to create
// a P2SH address, you pass data as a script.
//
//	 import (
//		  "github.com/kklash/bitcoinlib/address"
//		  "github.com/kklash/bitcoinlib/script"
//		  "encoding/hex"
//		  "log"
//	 )
//
//
//	 func main() {
//		  pk1, _ := hex.DecodeString("023bab8446121c93bca9f38d58d59fb399276e620ac760124bcf5b62a16f4e7e37")
//		  pk2, _ := hex.DecodeString("0390c351adb0dd332ff7eaf1456294fdeb75ee8f0b9e9a61bc860692e3698974d7")
//		  redeemScript := script.MakeP2MS(2, pk1, pk2)
//		  multisigAddress, err := address.Make(address.P2SH, redeemScript)
//		  if err != nil {
//		    log.Fatalln(err)
//		  }
//
//	  log.Println(multisigAddress)
//	 }
func Make(addressFormat constants.AddressFormat, data []byte) (string, error) {
	switch addressFormat {
	case constants.FormatP2PKH:
		return MakeP2PKHFromPublicKey(data)
	case constants.FormatP2SH:
		return MakeP2SHFromScript(data), nil
	case constants.FormatP2WPKH:
		return MakeP2WPKHFromPublicKey(data)
	case constants.FormatP2WSH:
		return MakeP2WSHFromScript(data)
	default:
		return "", ErrInvalidAddressFormat
	}
}

// MakeFromHash creates an address of the given addressFormat using a public-key or script-hash.
// For P2SH, P2PKH and P2WPKH address formats, hashed must be length 20. For P2WSH,
// hashed must be length 32. Returns ErrInvalidHash if the slice length does not match.
//
//	 import (
//		  "github.com/kklash/bitcoinlib/address"
//		  "github.com/kklash/bitcoinlib/constants"
//		  "github.com/kklash/bitcoinlib/bhash"
//		  "encoding/hex"
//		  "log"
//	 )
//
//	 func main() {
//		  pkHash, _ := hex.DecodeString("e79a6783b1546a18f7ef12c949d751b189cbf0e3")
//		  addrStr, err := address.MakeFromHash(constants.FormatP2PKH, pkHash)
//		  if err != nil {
//			  log.Fatalln(err)
//		  }
//
//	    log.Println(addrStr)
//	 }
func MakeFromHash(addressFormat constants.AddressFormat, hashed []byte) (string, error) {
	if hashed == nil {
		return "", ErrInvalidHash
	}

	var desiredHashLength int
	if addressFormat == constants.FormatP2WSH {
		desiredHashLength = 32
	} else {
		desiredHashLength = 20
	}
	if len(hashed) != desiredHashLength {
		return "", ErrInvalidHash
	}

	switch addressFormat {
	case constants.FormatP2PKH:
		var h [20]byte
		copy(h[:], hashed)
		return MakeP2PKHFromHash(h), nil
	case constants.FormatP2SH:
		var h [20]byte
		copy(h[:], hashed)
		return MakeP2SHFromHash(h), nil
	case constants.FormatP2WPKH:
		var h [20]byte
		copy(h[:], hashed)
		return MakeP2WPKHFromHash(h)
	case constants.FormatP2WSH:
		var h [32]byte
		copy(h[:], hashed)
		return MakeP2WSHFromHash(h)
	default:
		return "", ErrInvalidAddressFormat
	}
}

// DecodeBase58Address returns the hash and version number
// contained within a Base58 encoded address. Returns
// ErrInvalidAddress if the payload is not of the expected length.
func DecodeBase58Address(address string) (uint16, [20]byte, error) {
	var hashed [20]byte

	payload, err := base58check.Decode(address)
	if err != nil {
		return 0, hashed, err
	} else if len(payload) != 21 && len(payload) != 22 {
		return 0, hashed, ErrInvalidAddress
	}

	version := uint16(payload[0])
	payload = payload[1:]

	// Could be 2 byte version number
	if len(payload) == 21 {
		version <<= 8
		version |= uint16(payload[0])
		payload = payload[1:]
	}

	copy(hashed[:], payload)

	return version, hashed, nil
}

// DecodeBech32Address returns the prefix, version, and hash
// contained within a bech32 encoded address. Returns
// ErrInvalidAddress if the payload is not of the expected length.
func DecodeBech32Address(address string) (hrp string, version byte, payload []byte, err error) {
	hrp, version, payload, err = bech32.Decode(address)
	if err != nil {
		return
	}

	if len(payload) != 32 && len(payload) != 20 {
		err = ErrInvalidAddress
		return
	}

	return
}

// Decode validates and decodes an address to determine its format and script pub
// key (output locking script), by checking its embedded version numbers against
// the constants.CurrentNetwork. Returns the address format, the locking script
// which the address encodes, or and an error if the address is not valid.
func Decode(address string) (constants.AddressFormat, []byte, error) {
	b58Version, payload, err := DecodeBase58Address(address)
	if err == nil {
		switch b58Version {
		case constants.CurrentNetwork.ScriptHash:
			return constants.FormatP2SH, script.MakeP2SHFromHash(payload), nil
		case constants.CurrentNetwork.PubkeyHash:
			return constants.FormatP2PKH, script.MakeP2PKHFromHash(payload), nil
		default:
			err = fmt.Errorf("%w: unknown address version byte '0x%.2X'", ErrInvalidAddress, b58Version)
			return constants.FormatNONSTANDARD, nil, err
		}
	}

	hrp, witnessVersion, witnessProgram, err := DecodeBech32Address(address)
	if err != nil {
		err = fmt.Errorf("%w: failed to decode address as base58 or bech32", ErrInvalidAddress)
	} else if hrp != constants.CurrentNetwork.Bech32 {
		err = fmt.Errorf("%w: unexpected bech32 prefix '%s'", ErrInvalidAddress, hrp)
	} else if witnessVersion != constants.WitnessVersionZero {
		err = fmt.Errorf("%w: unexpected witness version '0x%.2X'", ErrInvalidAddress, witnessVersion)
	}
	if err != nil {
		return constants.FormatNONSTANDARD, nil, err
	}

	switch len(witnessProgram) {
	case 20:
		var keyHash [20]byte
		copy(keyHash[:], witnessProgram)
		return constants.FormatP2WPKH, script.MakeP2WPKHFromHash(keyHash), nil
	case 32:
		var scriptHash [32]byte
		copy(scriptHash[:], witnessProgram)
		return constants.FormatP2WSH, script.MakeP2WSHFromHash(scriptHash), nil
	default:
		err := fmt.Errorf("%w: unexpected witness program length %d", ErrInvalidAddress, len(witnessProgram))
		return constants.FormatNONSTANDARD, nil, err
	}
}
