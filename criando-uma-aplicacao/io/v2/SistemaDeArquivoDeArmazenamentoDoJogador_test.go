package main

import (
	"strings"
	"testing"
)

func TestaArmazenamentoDeSistemaDeArquivo(t *testing.T) {

    t.Run("/liga de um leitor", func(t *testing.T) {
        bancoDeDados := strings.NewReader(`[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)

        armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

        recebido := armazenamento.PegaLiga()

        esperado := []Jogador{
            {"Cleo", 10},
            {"Chris", 33},
        }

        defineLiga(t, recebido, esperado)
    })

		// ler novamente
		recebido = armazenamento.PegaLiga()
		defineLiga(t, recebido, esperado)
	})

	t.Run("/pega pontuacao do  jogador", func(t *testing.T) {
		bancoDeDados := strings.NewReader(`[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)

		armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

		recebido := armazenamento.("Chris")
		esperado := 33
		definePontuacaoIgual(t, recebido, esperado)
	})
}

func definePontuacaoIgual(t *testing.T, recebido, esperado int) {
	t.Helper()
	if recebido != esperado {
        t.Errorf("recebido '%s' esperado '%s'", recebido, esperado)
    }
}
