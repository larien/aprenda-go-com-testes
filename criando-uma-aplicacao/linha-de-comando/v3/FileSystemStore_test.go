package poker

import (
	"io/ioutil"
	"os"
	"testing"
)

func criarArquivoTemporario(t *testing.T, dadosIniciais string) (*os.File, func()) {
	t.Helper()

	arquivoTemporario, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	arquivoTemporario.Write([]byte(dadosIniciais))

	removerArquivo := func() {
		os.Remove(arquivoTemporario.Nome())
	}

	return arquivoTemporario, removerArquivo
}

func TestArmazenaSistemaDeArquivo(t *testing.T) {

	t.Run("league sorted", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

		verificaSemErro(t, err)

		obtido := armazenamento.ObterLiga()

		esperado := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		verificaLiga(t, obtido, esperado)

		// read again
		obtido = armazenamento.ObterLiga()
		verificaLiga(t, obtido, esperado)
	})

	t.Run("get player score", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

		verificaSemErro(t, err)

		obtido := armazenamento.ObterPontuacaoDeJogador("Chris")
		esperado := 33
		verificaPontuacaoIgual(t, obtido, esperado)
	})

	t.Run("armazenamento wins for existing players", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

		verificaSemErro(t, err)

		armazenamento.RecordWin("Chris")

		obtido := armazenamento.ObterPontuacaoDeJogador("Chris")
		esperado := 34
		verificaPontuacaoIgual(t, obtido, esperado)
	})

	t.Run("armazenamento wins for existing players", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

		verificaSemErro(t, err)

		armazenamento.RecordWin("Pepper")

		obtido := armazenamento.ObterPontuacaoDeJogador("Pepper")
		esperado := 1
		verificaPontuacaoIgual(t, obtido, esperado)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, "")
		defer limparBaseDeDados()

		_, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

		verificaSemErro(t, err)
	})
}

func verificaPontuacaoIgual(t *testing.T, obtido, esperado int) {
	t.Helper()
	if obtido != esperado {
		t.Errorf("obtido %d esperado %d", obtido, esperado)
	}
}

func verificaSemErro(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but obtido one, %v", err)
	}
}
