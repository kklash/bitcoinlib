package common

// ReverseBytes returns a reversed copy of a byte slice.
func ReverseBytes(b []byte) []byte {
	if b == nil {
		return nil
	}

	bLength := len(b)
	reversed := make([]byte, bLength)
	for i := 0; i < bLength; i++ {
		reversed[bLength-(i+1)] = b[i]
	}

	return reversed
}

// Reverse the byte-order of a byte slice.
func ReverseBytesInPlace(sliceToReverse []byte) {
	var (
		buf       byte
		fromEnd   int
		fromStart int
	)

	for fromStart = 0; fromStart < len(sliceToReverse)/2; fromStart++ {
		fromEnd = len(sliceToReverse) - 1 - fromStart
		buf = sliceToReverse[fromStart]
		sliceToReverse[fromStart] = sliceToReverse[fromEnd]
		sliceToReverse[fromEnd] = buf
	}
}
