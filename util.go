package main

import "errors"

var ErrInvalidHexChar = errors.New("terdapat karakter hex yang tidak valid")
var ErrOddHexStringLength = errors.New("string Hex harus berpanjang genap")
var ErrNumberCannotConvertToHex = errors.New("angka tidak dapat direpresentasikan sebagai karakter hex")
var ErrBase64Multiple4 = errors.New("ukuran input harus kelipatan 4")
var ErrBase64InvalidChar = errors.New("ditemukan karakter base64 tidak valid")
var ErrBase64InvalidPadding = errors.New("padding tidak valid")

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
	return GCD(key_a, m) == 1
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

const base64 string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var unbase64 = [132]int{
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 0-11
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 12-23
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 24-35
	-1, -1, -1, -1, -1, -1, -1, 62, -1, -1, -1, 63, // 36-47
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -2, // 48-59
	-1, 0, -1, -1, -1, 0, 1, 2, 3, 4, 5, 6, // 60-71
	7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, // 72-83
	19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, -1, // 84-95
	-1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, // 96-107
	37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, // 108-119
	49, 50, 51, -1, -1, -1, -1, -1, -1, -1, -1, -1, // 120-131
}

func encodeBase64(bytes []uint8) string {
	result_len := len(bytes) / 3
	if len(bytes)%3 != 0 {
		result_len++
	}
	result_len *= 4
	result := make([]uint8, result_len)

	j := 0
	for i := 0; i < len(bytes); i += 3 {
		// 11111100
		b := (bytes[i] & 0xFC) >> 2
		result[j] = base64[b]
		j++

		if i+1 == len(bytes) {
			// 00000011
			b = (bytes[i] & 0x3) << 4
			result[j] = base64[b]
			j++
			result[j] = '='
			j++
			result[j] = '='
			j++
			break
		}

		// 00000011 11110000
		b = ((bytes[i] & 0x3) << 4) | ((bytes[i+1] & 0xF0) >> 4)
		result[j] = base64[b]
		j++

		if i+2 == len(bytes) {
			// 00001111
			b = (bytes[i+1] & 0xF) << 2
			result[j] = base64[b]
			j++
			result[j] = '='
			j++
			break
		}

		// 00001111 11000000
		b = ((bytes[i+1] & 0xF) << 2) | ((bytes[i+2] & 0xC0) >> 6)
		result[j] = base64[b]
		j++

		// 00111111
		b = (bytes[i+2] & 0x3F)
		result[j] = base64[b]
		j++
	}

	return string(result)
}

func decodeBase64Char(char uint8) (uint8, error) {
	if int(char) >= len(unbase64) {
		return 0, ErrBase64InvalidChar
	}
	res := unbase64[char]
	if res == -1 {
		return 0, ErrBase64InvalidChar
	}
	return uint8(res), nil
}

func decodeBase64(input string) ([]uint8, error) {
	if len(input)%4 != 0 {
		return nil, ErrBase64Multiple4
	}

	result_len := 3 * (len(input) / 4)

	result := make([]uint8, result_len)

	j := 0
	for i := 0; i < len(input); i += 4 {
		a, err := decodeBase64Char(input[i])
		if err != nil {
			return nil, err
		}

		b, err := decodeBase64Char(input[i+1])
		if err != nil {
			return nil, err
		}

		c, err := decodeBase64Char(input[i+2])
		if err != nil {
			return nil, err
		}

		d, err := decodeBase64Char(input[i+3])
		if err != nil {
			return nil, err
		}

		// aaaaaabb
		// 00111111 00110000
		o := (a << 2) | ((b & 0x30) >> 4)
		result[j] = o
		j++

		if input[i+2] == '=' {
			if input[i+3] != '=' {
				return nil, ErrBase64InvalidPadding
			}
			break
		}

		// bbbbcccc
		// 00001111 00111100
		o = ((b & 0xF) << 4) | ((c & 0x3C) >> 2)
		result[j] = o
		j++

		if input[i+3] == '=' {
			break
		}

		// ccdddddd
		// 000000011 00111111
		o = ((c & 0x3) << 6) | d
		result[j] = o
		j++
	}

	return result[:j], nil
}
