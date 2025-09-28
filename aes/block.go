package aes

type aesBlock struct {
	d [4][4]uint8
}

const blockLength uint8 = 16

func bytesToBlock(bytes [16]uint8) aesBlock {
	return aesBlock{
		d: [4][4]uint8{
			{bytes[0], bytes[1], bytes[2], bytes[3]},
			{bytes[4], bytes[5], bytes[6], bytes[7]},
			{bytes[8], bytes[9], bytes[10], bytes[11]},
			{bytes[12], bytes[13], bytes[14], bytes[15]},
		},
	}
}

func (block *aesBlock) toBytes() [16]uint8 {
	return [16]uint8{
		block.d[0][0], block.d[0][1], block.d[0][2], block.d[0][3],
		block.d[1][0], block.d[1][1], block.d[1][2], block.d[1][3],
		block.d[2][0], block.d[2][1], block.d[2][2], block.d[2][3],
		block.d[3][0], block.d[3][1], block.d[3][2], block.d[3][3],
	}
}

func (block *aesBlock) addRoundKey(round uint8, roundKey aesRoundKey) {
	for i := range uint8(4) {
		for j := range uint8(4) {
			block.d[j][i] ^= roundKey[(round*4)+j][i]
		}
	}
}

func (block *aesBlock) subBytes() {
	for i := range 4 {
		for j := range 4 {
			block.d[j][i] = getSbox(block.d[j][i])
		}
	}
}

func (block *aesBlock) mixColumns() {
	M2 := gfieldMult2
	M3 := gfieldMult3

	for i := range 4 {
		a := block.d[i][0]
		b := block.d[i][1]
		c := block.d[i][2]
		d := block.d[i][3]

		block.d[i][0] = M2(a) ^ M3(b) ^ c ^ d
		block.d[i][1] = a ^ M2(b) ^ M3(c) ^ d
		block.d[i][2] = a ^ b ^ M2(c) ^ M3(d)
		block.d[i][3] = M3(a) ^ b ^ c ^ M2(d)
	}
}

func (block *aesBlock) getRowElements(i uint8) (uint8, uint8, uint8, uint8) {
	return block.d[0][i], block.d[1][i], block.d[2][i], block.d[3][i]
}

func (block *aesBlock) setRowElements(i uint8, data [4]uint8) {
	block.d[0][i] = data[0]
	block.d[1][i] = data[1]
	block.d[2][i] = data[2]
	block.d[3][i] = data[3]
}

func (block *aesBlock) shiftRows() {
	a, b, c, d := block.getRowElements(1)
	block.setRowElements(1, [4]uint8{b, c, d, a})

	a, b, c, d = block.getRowElements(2)
	block.setRowElements(2, [4]uint8{c, d, a, b})

	a, b, c, d = block.getRowElements(3)
	block.setRowElements(3, [4]uint8{d, a, b, c})
}

func (block *aesBlock) invShiftRows() {
	a, b, c, d := block.getRowElements(1)
	block.setRowElements(1, [4]uint8{d, a, b, c})

	a, b, c, d = block.getRowElements(2)
	block.setRowElements(2, [4]uint8{c, d, a, b})

	a, b, c, d = block.getRowElements(3)
	block.setRowElements(3, [4]uint8{b, c, d, a})
}

func (block *aesBlock) invSubBytes() {
	for i := range 4 {
		for j := range 4 {
			block.d[j][i] = getInvSbox(block.d[j][i])
		}
	}
}

func (block *aesBlock) invMixColumns() {
	M := gfieldMult

	for i := range 4 {
		a := block.d[i][0]
		b := block.d[i][1]
		c := block.d[i][2]
		d := block.d[i][3]

		block.d[i][0] = M(a, 14) ^ M(b, 11) ^ M(c, 13) ^ M(d, 9)
		block.d[i][1] = M(a, 9) ^ M(b, 14) ^ M(c, 11) ^ M(d, 13)
		block.d[i][2] = M(a, 13) ^ M(b, 9) ^ M(c, 14) ^ M(d, 11)
		block.d[i][3] = M(a, 11) ^ M(b, 13) ^ M(c, 9) ^ M(d, 14)
	}
}
