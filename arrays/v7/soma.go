package main

// Soma calcula o valor total dos números em um slice
func Soma(numeros []int) int {
	soma := 0
	for _, numero := range numeros {
		soma += numero
	}
	return soma
}

// SomaTodoOResto calcula as somas de todos menos o primeiro número
// da coleção de slices recebida
func SomaTodoOResto(numerosParaSomar ...[]int) []int {
	var somas []int
	for _, numeros := range numerosParaSomar {
		if len(numeros) == 0 {
			somas = append(somas, 0)
		} else {
			final := numeros[1:]
			somas = append(somas, Soma(final))
		}
	}

	return somas
}
