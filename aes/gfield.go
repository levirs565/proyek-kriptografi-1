package aes

var gfieldMultInvTable [256]uint8

func gfieldMult2(a uint8) uint8 {
	res := a << 1

	if (a>>7)&1 == 1 {
		res ^= 0x1b
	}

	return res
}

func gfieldMult3(a uint8) uint8 {
	return gfieldMult2(a) ^ a
}

func gfieldMult(a uint8, b uint8) uint8 {
	res := uint8(0)

	for i := 0; i < 8; i++ {
		if b&1 == 1 {
			res ^= a
		}

		a = gfieldMult2(a)

		b >>= 1
	}

	return res
}

func gfieldMultInvInit() {
	for i := 1; i <= 255; i++ {
		if gfieldMultInvTable[i] != 0 {
			continue
		}
		for j := 1; j <= 255; j++ {
			if gfieldMult(uint8(i), uint8(j)) == 1 {
				gfieldMultInvTable[i] = uint8(j)
				gfieldMultInvTable[j] = uint8(i)
				break
			}
		}
	}
}

func gfieldMultInv(a uint8) uint8 {
	return gfieldMultInvTable[a]
}
