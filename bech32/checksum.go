package bech32

import (
	"github.com/kklash/bitcoinlib/constants"
)

// def bech32_hrp_expand(s):
//   return [ord(x) >> 5 for x in s] + [0] + [ord(x) & 31 for x in s]

func bech32HrpExpand(hrp string) []uint5 {
	v := make([]uint5, 0)
	w := make([]uint5, 0)
	for _, c := range hrp {
		v = append(v, uint5(c>>5))
		w = append(w, uint5(c&31))
	}

	v = append(v, 0)
	v = append(v, w...)
	return v
}

// def bech32_polymod(values):
//   GEN = [0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3]
//   chk = 1
//   for v in values:
//     b = (chk >> 25)
//     chk = (chk & 0x1ffffff) << 5 ^ v
//     for i in range(5):
//       chk ^= GEN[i] if ((b >> i) & 1) else 0
//   return chk

func bech32Polymod(values []uint5) int {
	chk := 1
	for _, v := range values {
		b := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ int(v)
		for i := 0; i < len(constants.Bech32ChecksumGen); i++ {
			if (b>>i)&1 == 1 {
				chk ^= constants.Bech32ChecksumGen[i]
			} else {
				chk ^= 0
			}
		}
	}

	return chk
}

// def bech32_create_checksum(hrp, data):
//   values = bech32_hrp_expand(hrp) + data
//   polymod = bech32_polymod(values + [0,0,0,0,0,0]) ^ 1
//   return [(polymod >> 5 * (5 - i)) & 31 for i in range(6)]

func bech32CreateChecksum(hrp string, values []uint5) []uint5 {
	values = append(bech32HrpExpand(hrp), values...)
	polymod := bech32Polymod(append(values, 0, 0, 0, 0, 0, 0)) ^ 1

	checksum := make([]uint5, 6)
	for i := 0; i < len(checksum); i++ {
		checksum[i] = uint5((polymod >> (5 * (5 - i))) & 31)
	}

	return checksum
}

// def bech32_verify_checksum(hrp, data):
//   return bech32_polymod(bech32_hrp_expand(hrp) + data) == 1

func bech32VerifyChecksum(hrp string, values []uint5) bool {
	values = append(bech32HrpExpand(hrp), values...)
	return bech32Polymod(values) == 1
}
