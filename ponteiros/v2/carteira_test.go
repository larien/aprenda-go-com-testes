package main

import (
	"testing"
)

func TestCarteira(t *testing.T) {
	confirmaSaldo := func(t *testing.T, carteira Carteira, esperado Bitcoin) {
		t.Helper()
		resultado := carteira.Saldo()

		if resultado != esperado {
			t.Errorf("resultado %s, esperado %s", resultado, esperado)
		}
	}

	t.Run("Depositar", func(t *testing.T) {
		carteira := Carteira{}
		carteira.Depositar(Bitcoin(10))
		confirmaSaldo(t, carteira, Bitcoin(10))
	})

	t.Run("Retirar", func(t *testing.T) {
		carteira := Carteira{saldo: Bitcoin(20)}
		carteira.Retirar(10)
		confirmaSaldo(t, carteira, Bitcoin(10))
	})

}
