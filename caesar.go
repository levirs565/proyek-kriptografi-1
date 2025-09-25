package main

func SafeIntShift(data, shift, mod int) int {
	return ((data % mod) + mod + (shift % mod) + mod) % mod
}

func ShiftAlphaByte(data byte, shift int) byte {
	if data >= 'a' && data <= 'z' {
		return 'a' + byte(SafeIntShift(int(data-'a'), shift, 26))
	} else if data >= 'A' && data <= 'Z' {
		return 'A' + byte(SafeIntShift(int(data-'A'), shift, 26))
	}
	return data
}

func EncrpytBytes(data []byte, key int) []byte {
	result := make([]byte, len(data))

	for i, v := range data {
		result[i] = ShiftAlphaByte(v, key)
	}

	return result
}

func DecryptBytes(data []byte, key int) []byte {
	return EncrpytBytes(data, -key)
}
