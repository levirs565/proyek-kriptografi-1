package main

func safeIntShift(data, shift, mod int) int {
	return ((data % mod) + mod + (shift % mod) + mod) % mod
}

func isKoprima(key_a int) bool {
	return GCD(key_a , 26) == 1
}