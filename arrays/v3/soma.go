package main

// Soma calcula o valor total dos números em um slice
func Soma(numeros []int) int {
	soma := 0
	for _, numero := range numeros {
		soma += numero
	}
	return soma
}
