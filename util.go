package main

import "errors"

var ErrInvalidHexChar = errors.New("terdapat karakter hex yang tidak valid")
var ErrOddHexStringLength = errors.New("string Hex harus berpanjang genap")
var ErrNumberCannotConvertToHex = errors.New("angka tidak dapat direpresentasikan sebagai karakter hex")

func safeIntShift(data, shift, mod int) int {
	return ((data % mod) + mod + (shift % mod) + mod) % mod
}

func isKoprima(key_a int, affineOption string) bool {
	var m int
	switch affineOption {
	case "Alfabet (A-Z)":
		m = 26
	case "ASCII":
		m = 255
	}
	return GCD(key_a , m) == 1
}

func decodeHexChar(ch byte) (uint8, error) {
	if ch >= '0' && ch <= '9' {
		return uint8(ch - '0'), nil
	} else if ch >= 'a' && ch <= 'f' {
		return uint8(ch-'a') + 10, nil
	} else if ch >= 'A' && ch <= 'F' {
		return uint8(ch-'A') + 10, nil
	}
	return 0, ErrInvalidHexChar
}

func decodeHexString(hex string) ([]uint8, error) {
	if len(hex)%2 == 1 {
		return nil, ErrOddHexStringLength
	}

	result := make([]uint8, len(hex)/2)
	for i := 0; i < len(hex); i += 2 {
		a, err := decodeHexChar(hex[i])
		if err != nil {
			return nil, err
		}
		b, err := decodeHexChar(hex[i+1])
		if err != nil {
			return nil, err
		}

		result[i/2] = (a << 4) | b
	}

	return result, nil
}

func encodeHexChar(b uint8) (byte, error) {
	if b <= 9 {
		return '0' + b, nil
	} else if b <= 15 {
		return 'a' + (b - 10), nil
	}
	return 0, ErrNumberCannotConvertToHex
}

func encodeHexString(bytes []uint8) (string, error) {
	result := make([]byte, len(bytes)*2)

	for i := range len(bytes) {
		b := bytes[i]
		c1, err := encodeHexChar(b >> 4)
		if err != nil {
			return "", err
		}
		c2, err := encodeHexChar(b & 0x0f)
		if err != nil {
			return "", err
		}
		result[i*2] = c1
		result[i*2+1] = c2
	}

	return string(result), nil
}
