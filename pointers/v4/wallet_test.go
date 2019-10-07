package main

import (
	"testing"
)

func TestCarteira(t *testing.T) {

	t.Run("Depositar", func(t *testing.T) {
		carteira := Carteira{}
		carteira.Depositar(Bitcoin(10))

		confirmarSaldo(t, carteira, Bitcoin(10))
	})

	t.Run("Retirar saldo suficiente", func(t *testing.T) {
		carteira := Carteira{Bitcoin(20)}
		erro := carteira.Retirar(Bitcoin(10))

		confirmarSaldo(t, carteira, Bitcoin(10))
		confirmarErroInexistente(t, err)
	})

	t.Run("Retirar saldo insuficiente", func(t *testing.T) {
		saldoInicial  := Bitcoin(20)
		carteira := Carteira{saldoInicial }
		erro := carteira.Retirar(Bitcoin(100))

		confirmarSaldo(t, carteira, saldoInicial )
		confirmarErro(t, err, ErrInsufficientFunds)
	})
}

func confirmarSaldo(t *testing.T, carteira Carteira, valorEsperado Bitcoin) {
	t.Helper()
	valor := carteira.Balance()

	if valor!= valoresperado {
		t.Errorf("valor'%s' valorEsperado '%s'", got, want)
	}
}

func confirmarErroInexistente(t *testing.T, valorerror) {
	t.Helper()
	if valor!= nil {
		t.Fatal("Recebeu um erro inesperado")
	}
}

func confirmarErro(t *testing.T, valorerror, valorEsperado error) {
	t.Helper()
	if valor== nil {
		t.Fatal("Esperava um erro mas nenhum ocorreu")
	}

	if valor!= valorEsperado {
		t.Errorf("valor'%s', valorEsperado '%s'", got, want)
	}
}
