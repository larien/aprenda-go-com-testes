package main

import (
	"testing"
)

func TestCarteira(t *testing.T) {

	confirmarSaldo := func(t *testing.T, carteira Carteira, valorEsperado Bitcoin) {
		t.Helper()
		valor := carteira.Saldo()

		if valor != valorEsperado {
			t.Errorf("valor %s valorEsperado %s", valor, valorEsperado)
		}
	}

	confirmarErro := func(t *testing.T, erro error) {
		t.Helper()
		if erro == nil {
			t.Error("Esperava um erro mas nenhum ocorreu.")
		}
	}

	t.Run("Deposit", func(t *testing.T) {
		carteira := Carteira{}
		carteira.Depositar(Bitcoin(10))

		confirmarSaldo(t, carteira, Bitcoin(10))
	})

	t.Run("Retirar com saldo suficiente", func(t *testing.T) {
		carteira := Carteira{Bitcoin(20)}
		carteira.Retirar(Bitcoin(10))

		confirmarSaldo(t, carteira, Bitcoin(10))
	})

	t.Run("Retirar com saldo insuficiente", func(t *testing.T) {
		saldoInicial := Bitcoin(20)
		carteira := Carteira{saldoInicial}
		erro := carteira.Retirar(Bitcoin(100))

		confirmarSaldo(t, carteira, saldoInicial)
		confirmarErro(t, erro)
	})
}
