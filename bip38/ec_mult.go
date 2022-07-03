package bip38

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	"github.com/kklash/bitcoinlib/address"
	"github.com/kklash/bitcoinlib/base58check"
	"github.com/kklash/bitcoinlib/bhash"
	"github.com/kklash/bitcoinlib/common"
	"github.com/kklash/bitcoinlib/ecc"
	"golang.org/x/crypto/scrypt"
)

var (
	// base58 encoded as 'passphrase'
	intermediateCodeMagicBytesLotSequence = []byte{0x2C, 0xE9, 0xB3, 0xE1, 0xFF, 0x39, 0xE2, 0x51}
	intermediateCodeMagicBytes            = []byte{0x2C, 0xE9, 0xB3, 0xE1, 0xFF, 0x39, 0xE2, 0x53}
)

func encodeLotSequence(lot, sequence uint32) ([]byte, error) {
	if lot > 0xfffff {
		return nil, fmt.Errorf("invalid bip38 lot number")
	}
	if sequence > 0xfff {
		return nil, fmt.Errorf("invalid bip38 sequence number")
	}
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, (lot<<12)+sequence)
	return buf, nil
}

func GenerateIntermediateCodeWithLotSequence(random io.Reader, password string, lot, sequence uint32) (string, error) {
	lotSequence, err := encodeLotSequence(lot, sequence)
	if err != nil {
		return "", err
	}

	ownerSalt := make([]byte, 4)
	if _, err := io.ReadFull(random, ownerSalt); err != nil {
		return "", fmt.Errorf("failed to generate random data: %w", err)
	}

	prefactor, err := scrypt.Key([]byte(password), ownerSalt, 16384, 8, 8, 32)
	if err != nil {
		return "", err
	}

	ownerEntropy := concatBytes(ownerSalt, lotSequence)

	passFactor := bhash.DoubleSha256(concatBytes(prefactor, ownerEntropy))
	ppX, ppY := ecc.Curve.ScalarBaseMult(passFactor[:])
	passPoint := ecc.SerializePointCompressed(ppX, ppY)

	intermediateCode := concatBytes(
		intermediateCodeMagicBytesLotSequence,
		ownerEntropy,
		passPoint,
	)

	return base58check.Encode(intermediateCode), nil
}

func GenerateIntermediateCode(random io.Reader, password string) (string, error) {
	ownerEntropy := make([]byte, 8)
	if _, err := io.ReadFull(random, ownerEntropy); err != nil {
		return "", fmt.Errorf("failed to generate random data: %w", err)
	}

	passFactor, err := scrypt.Key([]byte(password), ownerEntropy, 16384, 8, 8, 32)
	if err != nil {
		return "", err
	}

	ppX, ppY := ecc.Curve.ScalarBaseMult(passFactor[:])
	passPoint := ecc.SerializePointCompressed(ppX, ppY)

	intermediateCode := concatBytes(
		intermediateCodeMagicBytes,
		ownerEntropy,
		passPoint,
	)

	return base58check.Encode(intermediateCode), nil
}

func EncryptIntermediateCode(random io.Reader, intermediateCodeStr string, compressed bool) (string, error) {
	intermediateCode, err := base58check.Decode(intermediateCodeStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode intermediate code string: %w", err)
	}

	if len(intermediateCode) != 49 {
		return "", fmt.Errorf("invalid intermediate code format")
	}

	payloadBuf := new(bytes.Buffer)
	payloadBuf.Write(prefixBytes(true))

	var useLotSequence bool
	if bytes.Equal(intermediateCode[:8], intermediateCodeMagicBytesLotSequence) {
		useLotSequence = true
	} else if bytes.Equal(intermediateCode[:8], intermediateCodeMagicBytes) {
		useLotSequence = false
	} else {
		return "", fmt.Errorf("expected intermediate code magic bytes prefix")
	}

	payloadBuf.WriteByte(
		encodeFlagByte(compressed, true, useLotSequence),
	)

	ownerEntropy := intermediateCode[8:16]
	passPoint := intermediateCode[16:]

	seedb := make([]byte, 24)
	if _, err := io.ReadFull(random, seedb); err != nil {
		return "", fmt.Errorf("failed to generate seedb: %w", err)
	}
	factorb := bhash.DoubleSha256(seedb)
	ppX, ppY, err := ecc.DeserializePoint(passPoint)
	if err != nil {
		return "", err
	}

	pubX, pubY := ecc.Curve.ScalarMult(ppX, ppY, factorb[:])
	publicKey := ecc.SerializePoint(pubX, pubY, compressed)

	p2pkhAddress, err := address.MakeP2PKHFromPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	addressHashFull := bhash.DoubleSha256([]byte(p2pkhAddress))
	addressHash := addressHashFull[:4]

	payloadBuf.Write(addressHash)
	payloadBuf.Write(ownerEntropy)

	scryptedKey, err := scrypt.Key(passPoint, concatBytes(addressHash, ownerEntropy), 1024, 1, 1, 64)
	if err != nil {
		return "", err
	}

	dk1, dk2 := scryptedKey[:32], scryptedKey[32:]

	block, err := aes.NewCipher(dk2)
	if err != nil {
		return "", err
	}

	encryptedHalf1 := make([]byte, 16)
	encryptedHalf2 := make([]byte, 16)

	block.Encrypt(encryptedHalf1, common.XorBytes(seedb[:16], dk1[:16]))

	v := concatBytes(encryptedHalf1[8:], seedb[16:])
	block.Encrypt(encryptedHalf2, common.XorBytes(v, dk1[16:]))

	payloadBuf.Write(encryptedHalf1[:8])
	payloadBuf.Write(encryptedHalf2)

	encryptedKeyString := base58check.Encode(payloadBuf.Bytes())
	return encryptedKeyString, nil
}

func decryptECMult(
	decodedEncryptedKey []byte,
	password string,
) (
	privateKey []byte,
	compressed bool,
	err error,
) {
	compressed = decodedEncryptedKey[2]&0b00100000 != 0
	useLotSequence := decodedEncryptedKey[2]&0b00000100 != 0
	addressHash := decodedEncryptedKey[3:7]
	ownerEntropy := decodedEncryptedKey[7:15]

	var passFactor []byte
	if useLotSequence {
		// lot and sequence
		ownerSalt := ownerEntropy[:4]
		prefactor, _ := scrypt.Key([]byte(password), ownerSalt, 16384, 8, 8, 32)
		_passFactor := bhash.DoubleSha256(concatBytes(prefactor, ownerEntropy))
		passFactor = _passFactor[:]
	} else {
		passFactor, _ = scrypt.Key([]byte(password), ownerEntropy, 16384, 8, 8, 32)
	}

	ppX, ppY := ecc.Curve.ScalarBaseMult(passFactor)
	passPoint := ecc.SerializePointCompressed(ppX, ppY)

	scryptedKey, err := scrypt.Key(passPoint, concatBytes(addressHash, ownerEntropy), 1024, 1, 1, 64)
	if err != nil {
		return nil, false, err
	}

	dk1, dk2 := scryptedKey[:32], scryptedKey[32:]

	block, err := aes.NewCipher(dk2)
	if err != nil {
		return nil, false, err
	}

	encryptedHalf1 := make([]byte, 16)
	copy(encryptedHalf1, decodedEncryptedKey[15:23])
	encryptedHalf2 := decodedEncryptedKey[23:]

	seedb := make([]byte, 24)

	decryptedHalf2 := make([]byte, 16)
	block.Decrypt(decryptedHalf2, encryptedHalf2)
	copy(encryptedHalf1[8:], common.XorBytes(decryptedHalf2[:8], dk1[16:24]))
	copy(seedb[16:], common.XorBytes(decryptedHalf2[8:], dk1[24:]))

	decryptedHalf1 := make([]byte, 16)
	block.Decrypt(decryptedHalf1, encryptedHalf1)
	copy(seedb[:16], common.XorBytes(decryptedHalf1, dk1[:16]))

	factorb := bhash.DoubleSha256(seedb)
	fb := new(big.Int).SetBytes(factorb[:])
	pf := new(big.Int).SetBytes(passFactor)
	fb.Mul(fb, pf)
	fb.Mod(fb, ecc.Curve.Params().N)
	privateKey = fb.FillBytes(make([]byte, 32))

	return privateKey, compressed, nil
}
