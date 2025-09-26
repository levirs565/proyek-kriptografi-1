package aes

type AESVariant int

const (
	AES128 AESVariant = iota
	AES192
	AES256
)

type AesContext struct {
	roundKey   [44][4]uint8
	iv         [16]uint8
	roundCount uint8
	keyLength  uint8
}

type aesBlock [][4]uint8
type aesRoundKey [][4]uint8

// [kolom][baris]

const nb uint8 = 4

var rcon [11]uint8 = [11]uint8{
	0x8d, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1b, 0x36,
}

func NewAesContext(variant AESVariant, key []uint8) AesContext {
	ctx := AesContext{}

	switch variant {
	case AES128:
		ctx.keyLength = 4
		ctx.roundCount = 10
	case AES192:
		ctx.keyLength = 6
		ctx.roundCount = 12
	default:
		ctx.keyLength = 8
		ctx.roundCount = 14
	}

	ctx.keyExpansion(key)
	return ctx
}

func (c *AesContext) keyExpansion(key []uint8) {
	i := uint8(0)
	for i = 0; i < c.keyLength; i++ {
		j := i * 4
		c.roundKey[i][0] = key[j+0]
		c.roundKey[i][1] = key[j+1]
		c.roundKey[i][2] = key[j+2]
		c.roundKey[i][3] = key[j+3]
	}

	var temp [4]uint8
	for i = c.keyLength; i < (nb * (c.roundCount + 1)); i++ {
		temp[0] = c.roundKey[i-1][0]
		temp[1] = c.roundKey[i-1][1]
		temp[2] = c.roundKey[i-1][2]
		temp[3] = c.roundKey[i-1][3]

		if i%c.keyLength == 0 {
			last := temp[0]
			temp[0] = temp[1]
			temp[1] = temp[2]
			temp[2] = temp[3]
			temp[3] = last

			temp[0] = getSbox(temp[0])
			temp[1] = getSbox(temp[1])
			temp[2] = getSbox(temp[2])
			temp[3] = getSbox(temp[3])

			temp[0] = temp[0] ^ rcon[i/c.keyLength]
		}

		k := (i - c.keyLength)
		c.roundKey[i][0] = c.roundKey[k][0] ^ temp[0]
		c.roundKey[i][1] = c.roundKey[k][1] ^ temp[1]
		c.roundKey[i][2] = c.roundKey[k][2] ^ temp[2]
		c.roundKey[i][3] = c.roundKey[k][3] ^ temp[3]
	}
}

func addRoundKey(round uint8, block aesBlock, roundKey aesRoundKey) {
	for i := uint8(0); i < 4; i++ {
		for j := uint8(0); j < 4; j++ {
			block[j][i] ^= roundKey[(round*4)+j][i]
		}
	}
}

func subBytes(block aesBlock) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			block[j][i] = getSbox(block[j][i])
		}
	}
}

func shiftRows(block aesBlock) {
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

func mixColumns(block aesBlock) {
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

func (c *AesContext) cipher(block aesBlock, roundKey aesRoundKey) {
	round := uint8(0)
	addRoundKey(round, block, roundKey)

	for round = 1; ; round++ {
		subBytes(block)
		shiftRows(block)
		if round == c.roundCount {
			break
		}
		mixColumns(block)
		addRoundKey(round, block, roundKey)
	}

	addRoundKey(c.roundCount, block, roundKey)
}

func invShiftRows(block aesBlock) {
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

func invSubBytes(block aesBlock) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			block[j][i] = getInvSbox(block[j][i])
		}
	}
}

func invMixColumns(block aesBlock) {
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

func (c *AesContext) invCipher(block aesBlock, roundKey aesRoundKey) {
	round := uint8(c.roundCount)

	addRoundKey(round, block, roundKey)

	for round := c.roundCount - 1; ; round-- {
		invShiftRows(block)
		invSubBytes(block)
		addRoundKey(round, block, roundKey)
		if round == 0 {
			break
		}
		invMixColumns(block)
	}
}

func (c *AesContext) EncryptECBBlock(bytes [16]uint8) [16]uint8 {
	block := aesBlock{
		{bytes[0], bytes[1], bytes[2], bytes[3]},
		{bytes[4], bytes[5], bytes[6], bytes[7]},
		{bytes[8], bytes[9], bytes[10], bytes[11]},
		{bytes[12], bytes[13], bytes[14], bytes[15]},
	}

	c.cipher(block, c.roundKey[:])

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

	c.invCipher(block, c.roundKey[:])

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
