package aes

type AesContext struct {
	roundKey [44][4]uint8
	iv       [16]uint8
}

type aesBlock [][4]uint8
type aesRoundKey [][4]uint8

// [kolom][baris]

var nb uint8 = 4
var nk uint8 = 4
var nr uint8 = 10

var rcon [11]uint8 = [11]uint8{
	0x8d, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1b, 0x36,
}

func keyExpansion(key []uint8, roundKey aesRoundKey) {
	i := uint8(0)
	for i = 0; i < nk; i++ {
		j := i * 4
		roundKey[i][0] = key[j+0]
		roundKey[i][1] = key[j+1]
		roundKey[i][2] = key[j+2]
		roundKey[i][3] = key[j+3]
	}

	var temp [4]uint8
	for i = nk; i < (nb * (nr + 1)); i++ {
		temp[0] = roundKey[i-1][0]
		temp[1] = roundKey[i-1][1]
		temp[2] = roundKey[i-1][2]
		temp[3] = roundKey[i-1][3]

		if i%nk == 0 {
			last := temp[0]
			temp[0] = temp[1]
			temp[1] = temp[2]
			temp[2] = temp[3]
			temp[3] = last

			temp[0] = getSbox(temp[0])
			temp[1] = getSbox(temp[1])
			temp[2] = getSbox(temp[2])
			temp[3] = getSbox(temp[3])

			temp[0] = temp[0] ^ rcon[i/nk]
		}

		k := (i - nk)
		roundKey[i][0] = roundKey[k][0] ^ temp[0]
		roundKey[i][1] = roundKey[k][1] ^ temp[1]
		roundKey[i][2] = roundKey[k][2] ^ temp[2]
		roundKey[i][3] = roundKey[k][3] ^ temp[3]
	}
}

func aesAddRoundKey(round uint8, block aesBlock, roundKey aesRoundKey) {
	for i := uint8(0); i < 4; i++ {
		for j := uint8(0); j < 4; j++ {
			block[j][i] ^= roundKey[(round*4)+j][i]
		}
	}
}

func aesSubBytes(block aesBlock) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			block[j][i] = getSbox(block[j][i])
		}
	}
}

func aesShiftRows(block aesBlock) {
	// a, b, c, d
	// b, c, d, a
	temp := block[0][1]
	block[0][1] = block[1][1]
	block[1][1] = block[2][1]
	block[2][1] = block[3][1]
	block[3][1] = temp

	// a, b, c, d
	// c, d, a, b
	temp = block[0][2]
	block[0][2] = block[2][2]
	block[2][2] = temp

	temp = block[1][2]
	block[1][2] = block[3][2]
	block[3][2] = temp

	// a, b, c, d
	// d, a, b, c
	temp = block[0][3]
	block[0][3] = block[3][3]
	block[3][3] = block[2][3]
	block[2][3] = block[1][3]
	block[1][3] = temp
}

func aesMixColumns(block aesBlock) {
	M2 := gfieldMult2
	M3 := gfieldMult3

	for i := 0; i < 4; i++ {
		a := block[i][0]
		b := block[i][1]
		c := block[i][2]
		d := block[i][3]

		block[i][0] = M2(a) ^ M3(b) ^ c ^ d
		block[i][1] = a ^ M2(b) ^ M3(c) ^ d
		block[i][2] = a ^ b ^ M2(c) ^ M3(d)
		block[i][3] = M3(a) ^ b ^ c ^ M2(d)
	}
}

func aesCipher(block aesBlock, roundKey aesRoundKey) {
	round := uint8(0)
	aesAddRoundKey(round, block, roundKey)

	for round = 1; ; round++ {
		aesSubBytes(block)
		aesShiftRows(block)
		if round == nr {
			break
		}
		aesMixColumns(block)
		aesAddRoundKey(round, block, roundKey)
	}

	aesAddRoundKey(nr, block, roundKey)
}

func aesInvShiftRows(block aesBlock) {
	// a, b, c, d
	// d, a, b, c
	temp := block[3][1]
	block[3][1] = block[2][1]
	block[2][1] = block[1][1]
	block[1][1] = block[0][1]
	block[0][1] = temp

	// a, b, c, d
	// c, d, a, b
	temp = block[0][2]
	block[0][2] = block[2][2]
	block[2][2] = temp

	temp = block[1][2]
	block[1][2] = block[3][2]
	block[3][2] = temp

	// a, b, c, d
	// b, c, d, a
	temp = block[0][3]
	block[0][3] = block[1][3]
	block[1][3] = block[2][3]
	block[2][3] = block[3][3]
	block[3][3] = temp
}

func aesInvSubBytes(block aesBlock) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			block[j][i] = getInvSbox(block[j][i])
		}
	}
}

func aesInvMixColumns(block aesBlock) {
	M := gfieldMult

	for i := 0; i < 4; i++ {
		a := block[i][0]
		b := block[i][1]
		c := block[i][2]
		d := block[i][3]

		block[i][0] = M(a, 14) ^ M(b, 11) ^ M(c, 13) ^ M(d, 9)
		block[i][1] = M(a, 9) ^ M(b, 14) ^ M(c, 11) ^ M(d, 13)
		block[i][2] = M(a, 13) ^ M(b, 9) ^ M(c, 14) ^ M(d, 11)
		block[i][3] = M(a, 11) ^ M(b, 13) ^ M(c, 9) ^ M(d, 14)
	}
}

func aesInvCipher(block aesBlock, roundKey aesRoundKey) {
	round := uint8(nr)

	aesAddRoundKey(round, block, roundKey)

	for round := nr - 1; ; round-- {
		aesInvShiftRows(block)
		aesInvSubBytes(block)
		aesAddRoundKey(round, block, roundKey)
		if round == 0 {
			break
		}
		aesInvMixColumns(block)
	}
}
func NewAesContext(key []uint8) AesContext {
	ctx := AesContext{}
	keyExpansion(key, ctx.roundKey[:])
	return ctx
}

func (c *AesContext) EncryptECBBlock(bytes [16]uint8) [16]uint8 {
	block := aesBlock{
		{bytes[0], bytes[1], bytes[2], bytes[3]},
		{bytes[4], bytes[5], bytes[6], bytes[7]},
		{bytes[8], bytes[9], bytes[10], bytes[11]},
		{bytes[12], bytes[13], bytes[14], bytes[15]},
	}

	aesCipher(block, c.roundKey[:])

	return [16]uint8{
		block[0][0], block[0][1], block[0][2], block[0][3],
		block[1][0], block[1][1], block[1][2], block[1][3],
		block[2][0], block[2][1], block[2][2], block[2][3],
		block[3][0], block[3][1], block[3][2], block[3][3],
	}
}

func (c *AesContext) DeryptECBBlock(bytes [16]uint8) [16]uint8 {
	block := aesBlock{
		{bytes[0], bytes[1], bytes[2], bytes[3]},
		{bytes[4], bytes[5], bytes[6], bytes[7]},
		{bytes[8], bytes[9], bytes[10], bytes[11]},
		{bytes[12], bytes[13], bytes[14], bytes[15]},
	}

	aesInvCipher(block, c.roundKey[:])

	return [16]uint8{
		block[0][0], block[0][1], block[0][2], block[0][3],
		block[1][0], block[1][1], block[1][2], block[1][3],
		block[2][0], block[2][1], block[2][2], block[2][3],
		block[3][0], block[3][1], block[3][2], block[3][3],
	}
}

func AESInit() {
	gfieldMultInvInit()
	sboxInit()
}
