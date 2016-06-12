package hack

import "testing"

func doTestIntToBits(val int, numBits int, expected string, t *testing.T) {
	if actual := intToBits(val, numBits); actual != expected {
		t.Errorf("strToBits(%v, %v) = %v; want %v", val, numBits, actual, expected)
	}
}

func TestIntToBits(t *testing.T) {
	type Cases struct {
		val      int
		numBits  int
		expected string
	}

	cases := []Cases{
		{1, 4, "0001"},
		{8, 4, "1000"},
		{0x10, 4, "0000"},
		{0x18, 4, "1000"},
	}

	for _, c := range cases {
		doTestIntToBits(c.val, c.numBits, c.expected, t)
	}
}
