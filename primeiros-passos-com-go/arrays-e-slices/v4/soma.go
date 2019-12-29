package main

// Soma calcula o valor total dos números em um slice
func Soma(numeros []int) int {
	soma := 0
	for _, numero := range numeros {
		soma += numero
	}
	return soma
}

// SomaTudo calcula as respectivas somas de cada slice recebido
func SomaTudo(numerosParaSomar ...[]int) []int {
	quantidadeDeNumeros := len(numerosParaSomar)
	somas := make([]int, quantidadeDeNumeros)

	for i, numeros := range numerosParaSomar {
		somas[i] = Soma(numeros)
	}

	return somas
}
