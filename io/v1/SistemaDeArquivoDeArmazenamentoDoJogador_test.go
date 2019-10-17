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
		//ler novamente
        defineLiga(t, recebido, esperado)
    })
}
