// Package bip38 provides private key encryption and decryption using scrypt and AES.
package bip38

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/kklash/bitcoinlib/address"
	"github.com/kklash/bitcoinlib/base58check"
	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/bip32"
	"github.com/kklash/bitcoinlib/ecc"
	"github.com/kklash/ekliptic"
	"golang.org/x/crypto/scrypt"
)

const (
	scryptN = 16384
	scryptR = 8
	scryptP = 8
)

var (
	// ErrInvalidPrivateKey is returned if attempting to encrypt a private key which
	// is either nil, or of the wrong length.
	ErrInvalidPrivateKey = errors.New("invalid private key to encrypt")

	// ErrInvalidEncryptedKey is returned when decoding a BIP38 key fails.
	ErrInvalidEncryptedKey = errors.New("key is not in bip38 encrypted format")

	// ErrDecryptionFailed is returned when decrypting a BIP38 key, either when the key is not
	// encrypted properly, or if the password is incorrect.
	ErrDecryptionFailed = errors.New("failed to decrypt bip38 key")
)

var secp = new(ekliptic.Curve)

func prefixBytes(ecMultiply bool) []byte {
	if ecMultiply {
		return []byte{0x01, 0x43}
	}

	return []byte{0x01, 0x42}
}

func xorBytes(b1, b2 []byte) []byte {
	size := len(b1)
	if len(b2) != size {
		panic("attempting to xor byte slices of different lengths")
	}

	output := make([]byte, size)
	for i := 0; i < size; i++ {
		output[i] = b1[i] ^ b2[i]
	}

	return output
}

func concatBytes(slices ...[]byte) []byte {
	totalSize := 0
	for _, slice := range slices {
		totalSize += len(slice)
	}
	result := make([]byte, totalSize)
	i := 0
	for _, slice := range slices {
		i += copy(result[i:], slice)
	}
	return result
}

func encodeUint32(n uint32) []byte {
	buf := make([]byte, 4)
	for i := 3; n > 0; i-- {
		buf[i] = byte(n & 0xff)
		n >>= 8
	}
	return buf
}

func deriveKey(password string, salt []byte, length int) []byte {
	dk, _ := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, length)
	return dk
}

func deriveAddress(privateKey []byte, compressed bool) (string, error) {
	publicKey := bip32.Neuter(privateKey, compressed)

	p2pkhAddress, err := address.MakeP2PKHFromPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	return p2pkhAddress, nil
}

func encodeFlagByte(compressed, ecMultiply, lotAndSequence bool) byte {
	var flagByte byte

	// For non-EC-multiplied keys, the first two bits of the flag byte are 11. For EC-multiplied keys, they are 00.
	if !ecMultiply {
		flagByte |= 0b11000000
	}
	// indicates the key should be converted to a bitcoin address using the compressed public key format.
	if compressed {
		flagByte |= 0b00100000
	}
	// indicates whether a lot and sequence number are encoded into the first factor, and activates special behavior
	// for including them in the decryption process. This applies to EC-multiplied keys only. Must be 0 for non-EC-multiplied keys.
	if ecMultiply && lotAndSequence {
		flagByte |= 0b00000100
	}

	return flagByte
}

func encodeLotSequence(lot, sequence uint32) []byte {
	if lot > 0xfffff {
		panic("invalid bip38 lot number")
	}
	if sequence > 0xfff {
		panic("invalid bip38 sequence number")
	}
	return encodeUint32((lot << 12) + sequence)
}

func generateIntermediateCodeWithLotSequence(password string, lot, sequence uint32) string {
	lotSequence := encodeLotSequence(lot, sequence)

	var ownerSalt [4]byte
	rand.Read(ownerSalt[:])
	prefactor := deriveKey(password, ownerSalt[:], 32)
	ownerEntropy := concatBytes(ownerSalt[:], lotSequence)

	passFactor := bhash.DoubleSha256(concatBytes(prefactor, ownerEntropy))
	ppX, ppY := secp.ScalarBaseMult(passFactor[:])
	passPoint := ecc.SerializePointCompressed(ppX, ppY)

	intermediateCode := concatBytes(
		[]byte{0x2C, 0xE9, 0xB3, 0xE1, 0xFF, 0x39, 0xE2, 0x51}, // base58 encoded as 'passphrase'
		ownerEntropy,
		passPoint,
	)

	return base58check.Encode(intermediateCode)
}

func generateIntermediateCode(password string) string {
	var ownerEntropy [8]byte
	rand.Read(ownerEntropy[:])
	passFactor := deriveKey(password, ownerEntropy[:], 32)

	ppX, ppY := secp.ScalarBaseMult(passFactor[:])
	passPoint := ecc.SerializePointCompressed(ppX, ppY)

	intermediateCode := concatBytes(
		[]byte{0x2C, 0xE9, 0xB3, 0xE1, 0xFF, 0x39, 0xE2, 0x53}, // base58 encoded as 'passphrase'
		ownerEntropy[:],
		passPoint,
	)

	return base58check.Encode(intermediateCode)
}

func encrypt(
	privateKey []byte,
	password string,
	compressed bool,
	ecMultiply bool,
	lotAndSequence bool,
) (string, error) {
	if len(privateKey) != 32 {
		return "", ErrInvalidPrivateKey
	}

	payloadBuf := new(bytes.Buffer)

	payloadBuf.Write(prefixBytes(ecMultiply))
	payloadBuf.WriteByte(
		encodeFlagByte(compressed, ecMultiply, lotAndSequence),
	)

	expectedAddress, err := deriveAddress(privateKey, compressed)
	if err != nil {
		return "", err
	}

	hashedAddress := bhash.DoubleSha256([]byte(expectedAddress))
	hashedAddressSalt := hashedAddress[:4]

	payloadBuf.Write(hashedAddressSalt)

	scryptedKey := deriveKey(password, hashedAddressSalt, 64)
	dk1, dk2 := scryptedKey[:32], scryptedKey[32:]

	cipher, err := aes.NewCipher(dk2)
	if err != nil {
		return "", err
	}

	var (
		encryptedHalf1 [16]byte
		encryptedHalf2 [16]byte
	)
	cipher.Encrypt(encryptedHalf1[:], xorBytes(privateKey[:16], dk1[:16]))
	cipher.Encrypt(encryptedHalf2[:], xorBytes(privateKey[16:], dk1[16:]))

	payloadBuf.Write(encryptedHalf1[:])
	payloadBuf.Write(encryptedHalf2[:])

	encryptedKeyString := base58check.Encode(payloadBuf.Bytes())
	return encryptedKeyString, nil
}

func decrypt(
	encryptedKeyString string,
	password string,
) (
	privateKey []byte,
	compressed bool,
	err error,
) {
	decodedEncryptedKey, err := base58check.Decode(encryptedKeyString)
	if err != nil {
		return nil, false, err
	}

	if len(decodedEncryptedKey) != 39 {
		return nil, false, fmt.Errorf("%w: incorrect bip38 key payload size", ErrInvalidEncryptedKey)
	}

	if decodedEncryptedKey[0] != 1 {
		return nil, false, fmt.Errorf("%w: key prefix byte incorrect", ErrInvalidEncryptedKey)
	}

	if decodedEncryptedKey[1] == 0x43 {
		return nil, false, fmt.Errorf("%w: ecMultiply keys are not supported", ErrInvalidEncryptedKey)
	} else if decodedEncryptedKey[1] != 0x42 {
		return nil, false, fmt.Errorf("%w: ecMultiply byte incorrect", ErrInvalidEncryptedKey)
	}

	if decodedEncryptedKey[2]&0b11000000 == 0 {
		return nil, false, fmt.Errorf("%w: ecMultiply keys are not supported", ErrInvalidEncryptedKey)
	}
	compressed = decodedEncryptedKey[2]&0b00100000 != 0
	hashedAddressSalt := decodedEncryptedKey[3:7]

	scryptedKey := deriveKey(password, hashedAddressSalt, 64)
	dk1, dk2 := scryptedKey[:32], scryptedKey[32:]

	cipher, err := aes.NewCipher(dk2)
	if err != nil {
		return nil, false, err
	}

	var decryptedPayload [32]byte

	cipher.Decrypt(decryptedPayload[:16], decodedEncryptedKey[7:23])
	cipher.Decrypt(decryptedPayload[16:], decodedEncryptedKey[23:])

	privateKey = make([]byte, 32)

	copy(privateKey[:16], xorBytes(decryptedPayload[:16], dk1[:16]))
	copy(privateKey[16:], xorBytes(decryptedPayload[16:], dk1[16:]))

	addr, err := deriveAddress(privateKey, compressed)
	if err != nil {
		return nil, false, err
	}

	hashedAddress := bhash.DoubleSha256([]byte(addr))
	if !bytes.Equal(hashedAddress[:4], hashedAddressSalt) {
		return nil, false, fmt.Errorf("%w: derived P2PKH address does not match checksum", ErrDecryptionFailed)
	}

	return privateKey, compressed, nil
}
