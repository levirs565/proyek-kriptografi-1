package aes

import "errors"

type AESVariant int

const (
	AES128 AESVariant = iota
	AES192
	AES256
)

type AesContext struct {
	roundKey   [60][4]uint8
	iv         [16]uint8
	roundCount uint8
	keyLength  uint8
}

type aesRoundKey [][4]uint8

// [kolom][baris]

var ErrInvalidLength = errors.New("ukuran bytes tidak sesuai dengan ukuran blok")

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

func rotWord(data [4]uint8) [4]uint8 {
	return [4]uint8{data[1], data[2], data[3], data[0]}
}

func subWord(data [4]uint8) [4]uint8 {
	return [4]uint8{
		getSbox(data[0]),
		getSbox(data[1]),
		getSbox(data[2]),
		getSbox(data[3]),
	}
}

func xorWord(a, b [4]uint8) [4]uint8 {
	return [4]uint8{
		a[0] ^ b[0],
		a[1] ^ b[1],
		a[2] ^ b[2],
		a[3] ^ b[3],
	}
}

func (c *AesContext) keyExpansion(key []uint8) {
	for i := uint8(0); i < (nb * (c.roundCount + 1)); i++ {
		if i < c.keyLength {
			j := i * 4
			copy(c.roundKey[i][:], key[j:j+4])
		} else {
			current := c.roundKey[i-1]

			if i%c.keyLength == 0 {
				current = subWord(rotWord(current))
				current[0] ^= rcon[i/c.keyLength]
			} else if c.keyLength > 6 && i%c.keyLength == 4 {
				current = subWord(current)
			}

			current = xorWord(c.roundKey[i-c.keyLength], current)
			c.roundKey[i] = current
		}
	}
}

func (c *AesContext) cipher(block *aesBlock, roundKey aesRoundKey) {
	round := uint8(0)
	block.addRoundKey(round, roundKey)

	for round = 1; ; round++ {
		block.subBytes()
		block.shiftRows()
		if round == c.roundCount {
			break
		}
		block.mixColumns()
		block.addRoundKey(round, roundKey)
	}

	block.addRoundKey(c.roundCount, roundKey)
}

func (c *AesContext) invCipher(block *aesBlock, roundKey aesRoundKey) {
	round := uint8(c.roundCount)

	block.addRoundKey(round, roundKey)

	for round := c.roundCount - 1; ; round-- {
		block.invShiftRows()
		block.invSubBytes()
		block.addRoundKey(round, roundKey)
		if round == 0 {
			break
		}
		block.invMixColumns()
	}
}

func (c *AesContext) EncryptECBBlock(bytes [16]uint8) [16]uint8 {
	block := bytesToBlock(bytes)

	c.cipher(&block, c.roundKey[:])

	return block.toBytes()
}

func (c *AesContext) DeryptECBBlock(bytes [16]uint8) [16]uint8 {
	block := bytesToBlock(bytes)

	c.invCipher(&block, c.roundKey[:])

	return block.toBytes()
}

func (c *AesContext) EncryptECB(bytes []uint8) []uint8 {
	padded_bytes := pkcs7Padd(bytes)
	cipher_bytes := make([]uint8, len(padded_bytes))

	for i := 0; i < len(padded_bytes); i += int(blockLength) {
		plain_block := padded_bytes[i : i+int(blockLength)]
		cipher_block := c.EncryptECBBlock([16]uint8(plain_block))
		copy(cipher_bytes[i:i+int(blockLength)], cipher_block[:])
	}

	return cipher_bytes
}

func (c *AesContext) DecryptECB(bytes []uint8) ([]uint8, error) {
	if len(bytes)%int(blockLength) != 0 {
		return nil, ErrInvalidLength
	}

	plain_bytes := make([]uint8, len(bytes))

	for i := 0; i < len(plain_bytes); i += int(blockLength) {
		cipher_block := bytes[i : i+int(blockLength)]
		plain_block := c.DeryptECBBlock([16]uint8(cipher_block))
		copy(plain_bytes[i:i+int(blockLength)], plain_block[:])
	}

	return pkcs7Unpadd(plain_bytes)
}

func AESInit() {
	gfieldMultInvInit()
	sboxInit()
}
