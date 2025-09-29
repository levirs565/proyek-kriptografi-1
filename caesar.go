package main

type CaesarMode int

const (
	CaesarModeAlphabet = iota
	CaesarModeAlphanum
	CaesarModeASCII
	CaesarModeCustom
)

func caesarShiftAlphaByte(data byte, shift int) byte { //fungsi menggeser alfabet(a-z)
	if data >= 'a' && data <= 'z' {
		return 'a' + byte(safeIntShift(int(data-'a'), shift, 26)) // 'a' itu valuenya 97
	} else if data >= 'A' && data <= 'Z' {
		return 'A' + byte(safeIntShift(int(data-'A'), shift, 26)) // 'A' itu valuenya 65
	}
	return data
}

func caesarShiftAlphaNumByte(data byte, shift int) byte { //fungsi menggeser alfabet(a-z) dan angka (0-9)
	if data >= 'a' && data <= 'z' {
		return 'a' + byte(safeIntShift(int(data-'a'), shift, 26))
	} else if data >= 'A' && data <= 'Z' {
		return 'A' + byte(safeIntShift(int(data-'A'), shift, 26))
	} else if data >= '0' && data <= '9' {
		return '0' + byte(safeIntShift(int(data-'0'), shift, 10))
	}
	return data
}

func caesarShiftAsciiByte(data byte, shift int) byte { //fungsi menggeser karakter ascii
	return byte(safeIntShift(int(data), shift, 256))
}

func caesarShiftCustomByte(data byte, shift int, charset string) byte {
	index := -1
	for i := 0; i < len(charset); i++ {
		if data == charset[i] {
			index = i
			break
		}
	}

	if index == -1 {
		return data
	}

	newIndex := safeIntShift(index, shift, len(charset))
	return charset[newIndex]
}

func caesarEncryptBytes(data []byte, key int, mode CaesarMode, customCharset string) []byte {
	result := make([]byte, len(data))
	for i, v := range data {
		switch mode {
		case CaesarModeAlphabet:
			result[i] = caesarShiftAlphaByte(v, key)
		case CaesarModeAlphanum:
			result[i] = caesarShiftAlphaNumByte(v, key)
		case CaesarModeASCII:
			result[i] = caesarShiftAsciiByte(v, key)
		case CaesarModeCustom:
			result[i] = caesarShiftCustomByte(v, key, customCharset)
		default:
			result[i] = v
		}
	}
	return result
}

func caesarDecryptBytes(data []byte, key int, mode CaesarMode, customCharset string) []byte {
	return caesarEncryptBytes(data, -key, mode, customCharset)
}
