package aes

import (
	"errors"
	"slices"
)

type Padding int

const (
	PKCS7Padding Padding = iota
	NoPadding
)

var ErrNoPaddingSizeNotMatch = errors.New("ukuran buffer tidak sesuai dengan ukuran blok (16). gunakan PKCS7 atau sesauikan ukuran")
var ErrInvalidPadding = errors.New("ukuran padding tidak valid")
var ErrPaddingNotMatch = errors.New("padding tidak sesuai")

func pkcs7Padd(bytes []uint8) ([]uint8, error) {
	length := len(bytes)
	extra := int(blockLength) - (length % int(blockLength))
	result := make([]uint8, length+extra)
	copy(result, bytes)
	for i := range extra {
		result[length+i] = uint8(extra)
	}
	return result, nil
}

func pkcs7Unpadd(bytes []uint8) ([]uint8, error) {
	length := len(bytes)

	if length == 0 {
		return []uint8{}, nil
	}

	amount := bytes[length-1]

	if amount == 0 || int(amount) > length {
		return nil, ErrInvalidPadding
	}

	for i := range amount {
		if bytes[length-1-int(i)] != amount {
			return nil, ErrPaddingNotMatch
		}
	}

	return slices.Clone(bytes[:length-int(amount)]), nil
}

func padd(bytes []uint8, padding Padding) ([]uint8, error) {
	if padding == PKCS7Padding {
		return pkcs7Padd(bytes)
	}

	if len(bytes)%int(blockLength) != 0 {
		return nil, ErrNoPaddingSizeNotMatch
	}
	return slices.Clone(bytes), nil
}

func unpadd(bytes []uint8, padding Padding) ([]uint8, error) {
	if padding == PKCS7Padding {
		return pkcs7Unpadd(bytes)
	}

	if len(bytes)%int(blockLength) != 0 {
		return nil, ErrInvalidPadding
	}
	return slices.Clone(bytes), nil
}
