package aes

var sboxTable, sboxInvTable [256]uint8

func sboxAffineTransform(a uint8) uint8 {
	return a ^ (a << 1) ^ (a << 2) ^ (a << 3) ^ (a << 4) ^ (a >> 7) ^ (a >> 6) ^ (a >> 5) ^ (a >> 4) ^ 0x63
}

func sboxInit() {
	for i := 0; i <= 255; i++ {
		inv := gfieldMultInv(uint8(i))
		affine := sboxAffineTransform(inv)

		sboxTable[i] = affine
		sboxInvTable[affine] = uint8(i)
	}
}

func getSbox(a uint8) uint8 {
	return sboxTable[a]
}

func getInvSbox(a uint8) uint8 {
	return sboxInvTable[a]
}
