package main

import "testing"

func TestSum(t *testing.T) {
	t.Run("coleção de qualquer tamanho", func(t *testing.T) {

		numeros := []int{1, 2, 3}

		resultado := Soma(numeros)
		esperado := 6

		if resultado != esperado {
			t.Errorf("resultado %d, esperado %d, dado, %v", resultado, esperado, numeros)
		}
	})
}
