package main

func safeIntShift(data, shift, mod int) int {
	return ((data % mod) + mod + (shift % mod) + mod) % mod
}

func isKoprima(key_a int, affineOption string) bool {
	var m int
	switch affineOption {
	case "Alfabet (A-Z)":
		m = 26
	case "ASCII":
		m = 255
	}
	return GCD(key_a , m) == 1
}