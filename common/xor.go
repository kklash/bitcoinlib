package common

// XorBytes returns the bitwise XOR of the two given given byte slices.
func XorBytes(b1, b2 []byte) []byte {
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
