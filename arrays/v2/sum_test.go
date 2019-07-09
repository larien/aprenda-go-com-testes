package main

import "testing"

func TestSoma(t *testing.T) {

	numeros := [5]int{1, 2, 3, 4, 5}

	resultado := Soma(numeros)
	esperado := 15

	if resultado != esperado {
		t.Errorf("resultado %d, esperado %d, dado, %v", resultado, esperado, numeros)
	}
}
