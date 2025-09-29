package main

type AffineMode int

const (
	AffineModeAlphabet = iota
	AffineModeASCII
)

func affineGetModulo(mode AffineMode) int {
	var m int
	switch mode {
	case AffineModeAlphabet:
		m = 26
	case AffineModeASCII:
		m = 256
	}
	return m
}

func affineIsCoprime(keyA int, mode AffineMode) bool {
	m := affineGetModulo(mode)
	return GCD(keyA, m) == 1
}

func affineEncryptChar(x int, keyA int, keyB int, mode AffineMode) int {
	m := affineGetModulo(mode)
	y := (keyA*x + keyB) % m
	return y
}

func affineGetCharIndex(b byte, mode AffineMode) (x int, startByte byte) {
	startByte = 0
	switch mode {
	case AffineModeAlphabet:
		if b >= 'a' && b <= 'z' {
			startByte = 'a'
			x = int(b - 'a')
		} else if b >= 'A' && b <= 'Z' {
			startByte = 'A'
			x = int(b - 'A')
		} else {
			x = -1
		}
	case AffineModeASCII:
		x = int(b)
	}
	return
}

func affineEncryptBytes(data []byte, keyA int, keyB int, mode AffineMode) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		x, startByte := affineGetCharIndex(b, mode)
		if x != -1 {
			y := affineEncryptChar(x, keyA, keyB, mode)
			result[i] = startByte + byte(y)
		} else {
			result[i] = data[i]
		}
	}
	return result
}

func affineDecryptChar(keyB int, y int, aInvers int, mode AffineMode) int {
	mod := affineGetModulo(mode)
	x := (aInvers * ((y - keyB + mod) % mod)) % mod
	return x
}

func affineDecryptBytes(data []byte, keyA int, keyB int, mode AffineMode) []byte {
	result := make([]byte, len(data))
	aInvers := affineGetInvers(keyA, mode)
	for i, b := range data {
		y, startByte := affineGetCharIndex(b, mode)
		if y != -1 {
			x := affineDecryptChar(keyB, y, aInvers, mode)
			result[i] = startByte + byte(x)
		} else {
			result[i] = data[i]
		}
	}
	return result
}

func affineGetInvers(keyA int, mode AffineMode) int {
	m := affineGetModulo(mode)
	for x := 1; x < m; x++ {
		if (keyA*x)%m == 1 {
			return x
		}
	}
	return -1
}
