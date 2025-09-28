package aes

import "slices"

func pkcs7Padd(bytes []uint8) []uint8 {
	length := len(bytes)
	extra := int(blockLength) - (length % int(blockLength))
	result := make([]uint8, length+extra)
	copy(result, bytes)
	for i := range extra {
		result[length+i] = uint8(extra)
	}
	return result
}

func pkcs7Unpadd(bytes []uint8) ([]uint8, bool) {
	length := len(bytes)

	if length == 0 {
		return nil, false
	}

	amount := bytes[length-1]

	if amount == 0 || int(amount) > length {
		return nil, false
	}

	for i := range amount {
		if bytes[length-1-int(i)] != amount {
			return nil, false
		}
	}

	return slices.Clone(bytes[:length-int(amount)]), true
}
