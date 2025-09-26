package main

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
	return byte(safeIntShift(int(data), shift, 128))
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

func caesarEncryptBytes(data []byte, key int, caesarOption string, customCharset string) []byte {
	result := make([]byte, len(data))
	for i, v := range data {
		switch caesarOption {
		case "Alfabet (A-Z)":
			result[i] = caesarShiftAlphaByte(v, key)
		case "Alphanum (A-Z dan 0-9)":
			result[i] = caesarShiftAlphaNumByte(v, key)
		case "ASCII":
			result[i] = caesarShiftAsciiByte(v, key)
		case "Custom karakter":
			result[i] = caesarShiftCustomByte(v, key, customCharset)
		default:
			result[i] = v
		}
	}
	return result
}

func caesarDecryptBytes(data []byte, key int, caesarOption string, customCharset string) []byte {
	return caesarEncryptBytes(data, -key, caesarOption, customCharset)
}