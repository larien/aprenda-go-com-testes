package main

import (
	"testing"
)

func TestCateira(t *testing.T) {

	confirmarSaldo := func(t *testing.T, carteira Carteira, valorEsperado Bitcoin) {
		t.Helper()
		valor := carteira.Saldo()

		if valor != valorEsperado {
			t.Errorf("valor %s valorEsperado %s", valor, valorEsperado)
		}
	}

	t.Run("Depositar", func(t *testing.T) {
		carteira := Carteira{}
		carteira.Depositar(Bitcoin(10))
		confirmarSaldo(t, carteira, Bitcoin(10))
	})

	t.Run("Retirar", func(t *testing.T) {
		carteira := Carteira{saldo: Bitcoin(20)}
		carteira.Retirar(10)
		confirmarSaldo(t, carteira, Bitcoin(10))
	})

}
