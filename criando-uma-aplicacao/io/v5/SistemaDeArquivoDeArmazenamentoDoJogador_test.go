package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func criaArquivoTemporario(t *testing.T, dadoInicial string) (io.ReadWriteSeeker, func()) {
    t.Helper()

   arquivotmp, err := ioutil.TempFile("", "db")

    if err != nil {
        t.Fatalf("não foi possivel escrever o arquivo temporário %v", err)
    }

    arquivotmp.Write([]byte(dadoInicial))

    removeArquivo := func() {
        arquivotmp.Close()
        os.Remove(arquivotmp.Name())
    }

    return arquivotmp, removeArquivo
}

func TestaArmazenamentoDeSistemaDeArquivo(t *testing.T) {

    t.Run("liga de um leitor", func(t *testing.T) {
        bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)
        defer limpaBancoDeDados()

        armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

        recebido := armazenamento.PegaLiga()

        esperado := []Jogador{
            {"Cleo", 10},
            {"Chris", 33},
        }

        defineLiga(t, recebido, esperado)

        // ler novamente
        recebido = armazenamento.PegaLiga()
        defineLiga(t, recebido, esperado)
    })

	t.Run("retorna pontuação do jogador", func(t *testing.T) {
        bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)
        defer limpaBancoDeDados()

        armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

        recebido := armazenamento.PegaPontuacaoDoJogador("Chris")
        esperado := 33
        definePontuacaoIgual(t, recebido, esperado)
    })

	t.Run("armazena vitorias de jogadores existentes", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)
        defer limpaBancoDeDados()

		armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

		armazenamento.SalvaVitoria("Chris")

		recebido := armazenamento.PegaPontuacaoDoJogador("Chris")
        esperado := 34
        definePontuacaoIgual(t, recebido, esperado)
	})

	t.Run("armazena vitorias para novos jogadores", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)
        defer limpaBancoDeDados()

		armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

		armazenamento.SalvaVitoria("Pepper")

		recebido := armazenamento.PegaPontuacaoDoJogador("Pepper")
		esperado := 1
		definePontuacaoIgual(t, recebido, esperado)
	})
}

func definePontuacaoIgual(t *testing.T, recebido, esperado int) {
	t.Helper()
	if recebido != esperado {
        t.Errorf("recebido '%s' esperado '%s'", recebido, esperado)
    }
}
