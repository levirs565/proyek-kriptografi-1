package main
// func convertToInteger(data byte) int {
//   if data >= 'a' && data <= 'z' {
// 		return int(data - 'a')
// 	} else if data >= 'A' && data <= 'Z' {

// 		return int(data - 'A')
// 	}
// 	return -1
// }

func enkripsiChar(x int, key_a int, key_b int) int {
//   y := (key_a * x + key_b) % 26
	y := (key_a * x + key_b) % 255
  return y 
}

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b 
	}
	return a
}

func affineEncryptBytes(data []byte, key_a int, key_b int) []byte {
  result := make([]byte, len(data))
  for i := 0; i < len(data); i++ {
	// println(data[i])
	// x := convertToInteger(data[i])
	x := int(data[i])
	// println("sudah di konvert", x)
	if x != -1 {
	  y:= enkripsiChar(x, key_a, key_b)
		println(y)
	//   if data[i] >='a' && data[i] <= 'z' {
	// 	  result[i] = byte(y + 'a')
	//   } else if data[i] >='A' && data[i] <= 'Z' {
	// 	result[i] = byte(y + 'A')
	//   }

		result[i] = byte(y);
	} else {
	  println("salah satu karakter bukan huruf")
	}
  }
  return result
}

func decryptChar(keyB int, y int, aInvers int) int{
	x := aInvers * (y - keyB) % 255
	return x
}

func affineDecryptBytes(data []byte, key_a int, key_b int) []byte {
	result := make([]byte, len(data))
	aInvers := getAInvers(key_a)
	println("didapatkan a invers", aInvers)
	for i := 0; i < len(data); i++ {
	//   y := convertToInteger(data[i])
		y := int(data[i])
	  if y != -1 {
		x:= decryptChar(key_b, y, aInvers)
		println("setelah di dekrip", x)
		// if data[i] >='a' && data[i] <= 'z' {
		// 	result[i] = byte(x + 'a')
		// 	println(result[i])
		// } else if data[i] >='A' && data[i] <= 'Z' {
		//   result[i] = byte(x + 'A')
		//   println(result[i])
		// }
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


