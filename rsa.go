package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

type RSAValues struct {
	E, N, D, P, Q, Phi *big.Int
}

var ErrRSAKeyOdd = errors.New("ukuran bit harus genap")
var ErrRSAInvalidParameter = errors.New("input tidak cukup untuk enkripsi (membutuhkan M, E, dan N)")

func RSAGenerateKeys(bits int) (*RSAValues, error) {
	if bits%2 != 0 {
		return nil, ErrRSAKeyOdd
	}

	e := big.NewInt(65537)
	one := big.NewInt(1)

	var p, q, n, phi, d *big.Int

	for {
		primeBits := bits / 2
		var err error
		p, err = rand.Prime(rand.Reader, primeBits)
		if err != nil {
			return nil, err
		}

		q, err = rand.Prime(rand.Reader, primeBits)
		if err != nil {
			return nil, err
		}

		n = new(big.Int).Mul(p, q)

		pMinus1 := new(big.Int).Sub(p, one)
		qMinus1 := new(big.Int).Sub(q, one)
		phi = new(big.Int).Mul(pMinus1, qMinus1)

		if new(big.Int).Mod(phi, e).Cmp(big.NewInt(0)) != 0 {
			d = new(big.Int).ModInverse(e, phi)
			if d != nil {
				break
			}
		}
	}

	vals := &RSAValues{
		P:   p,
		Q:   q,
		N:   n,
		Phi: phi,
		E:   e,
		D:   d,
	}

	return vals, nil
}

func (vals *RSAValues) CalculateMissingValues() error {
	if vals.P != nil && vals.Q != nil {
		if vals.N == nil {
			vals.N = new(big.Int).Mul(vals.P, vals.Q)
		}
		one := big.NewInt(1)
		pMinus1 := new(big.Int).Sub(vals.P, one)
		qMinus1 := new(big.Int).Sub(vals.Q, one)
		vals.Phi = new(big.Int).Mul(pMinus1, qMinus1)
	}
	if vals.E != nil && vals.Phi != nil && vals.D == nil {
		vals.D = new(big.Int).ModInverse(vals.E, vals.Phi)
		if vals.D == nil {
			return fmt.Errorf("gagal menghitung D. E dan PHI mungkin tidak relatif prima")
		}
	}
	return nil
}

func RSADecrypt(c, d, n *big.Int) (*big.Int, error) {
	if c == nil || d == nil || n == nil {
		return nil, fmt.Errorf("input tidak cukup untuk dekripsi (membutuhkan C, D, dan N)")
	}
	return new(big.Int).Exp(c, d, n), nil
}

func RSAEncrypt(m, e, n *big.Int) (*big.Int, error) {
	if m == nil || e == nil || n == nil {
		return nil, ErrRSAInvalidParameter
	}
	return new(big.Int).Exp(m, e, n), nil
}

func StringToBigInt(s string) *big.Int {
	return new(big.Int).SetBytes([]byte(s))
}

func BigIntToString(n *big.Int) string {
	if n == nil {
		return ""
	}
	return string(n.Bytes())
}
