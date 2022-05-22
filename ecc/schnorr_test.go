package ecc

import "testing"

func TestIsEven(t *testing.T) {
	if isEven(bigIntFromHex("1")) {
		t.Errorf("1 should not be even")
	}
	if !isEven(bigIntFromHex("2")) {
		t.Errorf("2 should be even")
	}
	if isEven(bigIntFromHex("3")) {
		t.Errorf("3 should not be even")
	}
	if isEven(bigIntFromHex("ab3")) {
		t.Errorf("0xab3 should not be even")
	}
}

func TestLiftX(t *testing.T) {
	test := func(xHex, yHex string) {
		x := bigIntFromHex(xHex)
		expectedY := bigIntFromHex(yHex)

		actualY := liftX(x)
		if !equal(actualY, expectedY) {
			t.Errorf("Expected to expand x value %s to y value %s; got %x", xHex, yHex, actualY.Bytes())
		}
	}

	test(
		"f722f5f1dfdb41951882ea19dcb520bcb5ad1f529da3b209eebfc1211674e39a",
		"07de0c7c747fea6a923018ebec44fd2e0b6ef82921fff1ae5fc5b3669698739e",
	)
	test(
		"BE5E1AFE40A8B4019213EC9D3B562456113197EB64BF43438170BD32DC29EAD2",
		"7B4786F12502E181CE71D75085493FCA2E586CCEEAAAED6A0E3C96607DFB55C2",
	)
}

func BenchmarkLiftX(b *testing.B) {
	n := bigIntFromHex("f722f5f1dfdb41951882ea19dcb520bcb5ad1f529da3b209eebfc1211674e39a")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		liftX(n)
	}
}
