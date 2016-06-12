package hack

import "fmt"

func intToBits(val int, numBits int) string {
	out := make([]byte, numBits)

	for i := numBits - 1; i >= 0; i-- {
		if val&1 != 0 {
			out[i] = '1'
		} else {
			out[i] = '0'
		}
		val = val >> 1
	}

	return string(out)
}

type CodeGen struct {
}

func (c *CodeGen) EmitAInstruction(addr uint16) {
	fmt.Println("0" + intToBits(int(addr), 15))
}

func (c *CodeGen) EmitCInstruction(dest, comp, jump int) {
	fmt.Printf("111%v%v%v\n",
		intToBits(comp, 7), intToBits(dest, 3), intToBits(jump, 3))
}
