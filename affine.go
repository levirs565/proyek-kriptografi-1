package main

func getM(affineOption string) int {
	var m int
	switch affineOption {
	case "Alfabet (A-Z)":
		m = 26
	case "Alphanum (A-Z dan 0-9)":
		m = 10 + 26
	case "ASCII":
		m = 256
	}
	return m
}

func isKoprima(key_a int, affineOption string) bool {
	m := getM(affineOption)
	return GCD(key_a , m) == 1
}

func enkripsiChar(x int, key_a int, key_b int, affineOption string) int {
	m := getM(affineOption)
	y := (key_a * x + key_b) % m
  	return y 
}

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b 
	}
	return a
}

func affineEncryptBytes(data []byte, key_a int, key_b int, affineOption string) ([]byte, error) {
	result := make([]byte, len(data))
	for i := range data {
		var x int
		switch affineOption {
			case "Alfabet (A-Z)":
				if isAlphabet(data[i]) {
					if data[i] >= 'a' && data[i] <= 'z'  {
						x = int(data[i] - 'a')
					} else {
						x = int(data[i] - 'A')
					}

				} else {
					x = -1
				}
			case "Alphanum (A-Z dan 0-9)":
				if isAlphabet(data[i]) {
					if data[i] >= 'a' && data[i] <= 'z'  {
						x = int(data[i] - 'a')
					} else {
						x = int(data[i] - 'A')
					}

				} else if isNumeric(data[i]) {
					x = int(data[i] - '0')
				} else {
					x = -1
				} 
				
			case "ASCII":
				x = int(data[i])
		}
	
	if x != -1 {
		
		y := enkripsiChar(x, key_a, key_b, affineOption)
		
		switch affineOption {
			case "Alfabet (A-Z)":
				if data[i] >= 'a' && data[i] <= 'z'  {
					result[i] = byte(y + 'a')
				} else {
					result[i] = byte(y + 'A')
				}
			case "Alphanum (A-Z dan 0-9)":
				if isAlphabet(data[i]) {
					if data[i] >= 'a' && data[i] <= 'z'  {
						result[i] = byte(y + 'a')
					} else {
						result[i] = byte(y + 'A')
					}

				} else if isNumeric(data[i]) {
					result[i] = byte(y + '0')
				} 
			case "ASCII":
				result[i] = byte(y);
			default:
		}
	
	} else {
		result[i] = data[i]
	}
}
  return result, nil
}

func decryptChar(keyB int, y int, aInvers int) int{
	x := aInvers * (y - keyB) % 255
	return x
}

func affineDecryptBytes(data []byte, key_a int, key_b int) []byte {
	result := make([]byte, len(data))
	aInvers := getAInvers(key_a)
	for i := 0; i < len(data); i++ {
		y := int(data[i])
	  if y != -1 {
		x:= decryptChar(key_b, y, aInvers)
		println("setelah di dekrip", x)
		result[i] = byte(x)
	  } else {
		result[i] = data[i];
	  }
	}
	return result
}


func getAInvers(keyA int)  int {
	m:= 255
	for x := 1; x < m; x++ {
		if(keyA * x) % m == 1{
			return x;
		}
	}
	return -1;
}


