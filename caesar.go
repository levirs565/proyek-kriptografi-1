package main

func caesarShiftAlphaByte(data byte, shift int) byte {
	if data >= 'a' && data <= 'z' {
		return 'a' + byte(safeIntShift(int(data-'a'), shift, 26))
	} else if data >= 'A' && data <= 'Z' {
		return 'A' + byte(safeIntShift(int(data-'A'), shift, 26))
	}
	return data
}

func caesarEncryptBytes(data []byte, key int) []byte {
	result := make([]byte, len(data))

	for i, v := range data {
		result[i] = caesarShiftAlphaByte(v, key)
	}

	return result
}

func caesarDecryptBytes(data []byte, key int) []byte {
	return caesarEncryptBytes(data, -key)
}
