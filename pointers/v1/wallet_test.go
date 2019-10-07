package main

import (
	"testing"
)

func TestCarteira(t *testing.T) {

	carteira := Carteira{}

	carteira.Depositar(Bitcoin(10))

	valor := carteira.Saldo()

	valorEsperado := Bitcoin(10)

	if valor != valorEsperado {
		t.Errorf("valor %s valorEsperado %s", valor, valorEsperado)
	}
}
