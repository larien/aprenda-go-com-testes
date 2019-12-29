package main

import (
	"testing"
)

func TestCarteira(t *testing.T) {
	carteira := Carteira{}

	carteira.Depositar(Bitcoin(10))

	resultado := carteira.Saldo()

	esperado := Bitcoin(10)

	if resultado != esperado {
		t.Errorf("resultado %s, esperado %s", resultado, esperado)
	}
}
