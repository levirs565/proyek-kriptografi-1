package main

import (
	"crypto/rand"
	"errors"
	"kriptografi1/aes"
	"math/big"
	"slices"
	"strings"
)

var ErrSuperInvalidKey = errors.New("kunci super tidak valid")
var ErrSuperInvalidCipher = errors.New("cipher tidak valid")
var ErrSuperWrongKey = errors.New("kunci yang digunakan tidak sesuai")

type SuperKey struct {
	public, private string
}

const SuperPrefixPublic = "SPUB"
const SuperPrefixPrivate = "SPRIV"

// Public Key: Prefix-N-E
// Private Key: Prefix-N-D

func SuperGenerateKey() (*SuperKey, error) {
	rsa_key, err := RSAGenerateKeys(2048)
	if err != nil {
		return nil, err
	}
	nBase64 := encodeBase64(rsa_key.N.Bytes())
	publicKey := SuperPrefixPublic + "-" + nBase64 + "-" + encodeBase64(rsa_key.E.Bytes())
	privateKey := SuperPrefixPrivate + "-" + nBase64 + "-" + encodeBase64(rsa_key.D.Bytes())
	return &SuperKey{
		public:  publicKey,
		private: privateKey,
	}, nil
}

func SuperDecodePublicKey(key string) (*RSAValues, error) {
	split := strings.Split(key, "-")
	if len(split) != 3 {
		return nil, ErrSuperInvalidKey
	}
	if split[0] != SuperPrefixPublic {
		return nil, ErrSuperInvalidKey
	}

	n_bytes, err := decodeBase64(split[1])
	if err != nil {
		return nil, errors.Join(ErrSuperInvalidKey, err)
	}
	n := new(big.Int).SetBytes(n_bytes)

	e_bytes, err := decodeBase64(split[2])
	if err != nil {
		return nil, errors.Join(ErrSuperInvalidKey, err)
	}
	e := new(big.Int).SetBytes(e_bytes)

	return &RSAValues{
		N: n,
		E: e,
	}, nil
}

func SuperDecodePrivateKey(key string) (*RSAValues, error) {
	split := strings.Split(key, "-")
	if len(split) != 3 {
		return nil, ErrSuperInvalidKey
	}
	if split[0] != SuperPrefixPrivate {
		return nil, ErrSuperInvalidKey
	}

	n_bytes, err := decodeBase64(split[1])
	if err != nil {
		return nil, errors.Join(ErrSuperInvalidKey, err)
	}
	n := new(big.Int).SetBytes(n_bytes)

	d_bytes, err := decodeBase64(split[2])
	if err != nil {
		return nil, errors.Join(ErrSuperInvalidKey, err)
	}
	d := new(big.Int).SetBytes(d_bytes)

	return &RSAValues{
		N: n,
		D: d,
	}, nil
}

const SuperEncryptSentinel = "SUPERENCRYPT"

// Structure Cipher
// [2048-bit enkripsi RSA][16 byte IV AES][cipher AES]
// Isi RSA: S[32 byte AES password][1 byte key affine a][1 byte ket affine b][1 byte key caesar]SUPERENCRYPT
// Total isi RSA = 35 + 13 byte = 48 byte
// 2048-bit = 256 byte
// Ukuran Minimal Cipher = 256 + 16 + 16 (pkcs7 padding) = 288

func SuperEncrypt(rsaKey *RSAValues, bytes []uint8) ([]uint8, error) {
	var aesKey [32]uint8

	_, err := rand.Read(aesKey[:])

	if err != nil {
		return nil, err
	}

	maxA := big.NewInt(128)
	kBig, err := rand.Int(rand.Reader, maxA)

	if err != nil {
		return nil, err
	}

	k := kBig.Int64()
	a := uint8(2*k + 1)

	var b [1]uint8
	_, err = rand.Read(b[:])

	if err != nil {
		return nil, err
	}

	var c [1]uint8
	_, err = rand.Read(c[:])

	if err != nil {
		return nil, err
	}

	var keys [48]uint8
	keys[0] = 'S'
	copy(keys[1:33], aesKey[:])
	keys[33] = a
	keys[34] = b[0]
	keys[35] = c[0]
	copy(keys[36:], []uint8(SuperEncryptSentinel))

	m := new(big.Int)
	m.SetBytes(keys[:])

	encryptedKeys, err := RSAEncrypt(m, rsaKey.E, rsaKey.N)
	var encryptedKeysBytes [256]uint8
	encryptedKeys.FillBytes(encryptedKeysBytes[:])

	if err != nil {
		return nil, err
	}

	var iv [16]uint8

	_, err = rand.Read(iv[:])

	if err != nil {
		return nil, err
	}

	aesCtx, err := aes.NewAesContext(aes.AES256, aesKey[:])

	if err != nil {
		return nil, err
	}

	aesCtx.SetIv(iv)

	encryptedBytes, err := aesCtx.EncryptECB(bytes, aes.PKCS7Padding)

	if err != nil {
		return nil, err
	}

	cipher := make([]uint8, 0)
	cipher = append(cipher, encryptedKeysBytes[:]...)
	cipher = append(cipher, iv[:]...)
	cipher = append(cipher, encryptedBytes...)

	return cipher, nil
}

func SuperDecrypt(rsaKey *RSAValues, bytes []uint8) ([]uint8, error) {
	if len(bytes) < 288 || len(bytes)%16 != 0 {
		return nil, ErrSuperInvalidCipher
	}

	rsaEncrypted := slices.Clone(bytes[:256])
	aesIv := slices.Clone(bytes[256:272])
	aesEncrypted := slices.Clone(bytes[272:])

	rsaEncryptedNum := new(big.Int)
	rsaEncryptedNum.SetBytes(rsaEncrypted)

	rsaPlainNum, err := RSADecrypt(rsaEncryptedNum, rsaKey.D, rsaKey.N)

	if err != nil {
		return nil, err
	}

	rsaPlain := []uint8(rsaPlainNum.Bytes())

	if len(rsaPlain) != 48 {
		return nil, ErrSuperWrongKey
	}

	if rsaPlain[0] != 'S' || !slices.Equal(rsaPlain[36:], []uint8(SuperEncryptSentinel)) {
		return nil, ErrSuperWrongKey
	}

	aesKey := rsaPlain[1:33]

	aesCtx, err := aes.NewAesContext(aes.AES256, aesKey)

	if err != nil {
		return nil, err
	}

	aesCtx.SetIv([16]uint8(aesIv))

	return aesCtx.DecryptECB(aesEncrypted, aes.PKCS7Padding)
}
