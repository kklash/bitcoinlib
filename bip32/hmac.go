package bip32

import (
	"crypto/hmac"
	"crypto/sha512"
)

func hmacSha512(key, data []byte) []byte {
	h := hmac.New(sha512.New, key)
	if _, err := h.Write(data); err != nil {
		panic("HMAC-SHA512 failed to write: " + err.Error())
	}

	return h.Sum(nil)[:]
}
