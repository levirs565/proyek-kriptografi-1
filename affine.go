package main

import "errors"

func enkripsiChar(x int, key_a int, key_b int, affineOption string) int {
	var m int
	switch affineOption {
	case "Alfabet (A-Z)":
		m = 26
	case "ASCII":
		m = 255
	}
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
	var x int
  for i := range data {
	switch affineOption {
	case "Alfabet (A-Z)":
		x = int(data[i] - 'A')
	case "ASCII":
		x = int(data[i])
	}
	
	if x != -1 {
		y:= enkripsiChar(x, key_a, key_b, affineOption)
		switch affineOption {
		case "Alfabet (A-Z)":
			result[i] = byte(y + 'A')
		case "ASCII":
			result[i] = byte(y);
		}
	} else {
		return nil, errors.New("salah input karakter")
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


