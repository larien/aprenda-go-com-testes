package poquer_test

import (
	"io/ioutil"
	"os"
	"testing"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/websockets/v2"
)

func criarArquivoTemporario(t *testing.T, dadosIniciais string) (*os.File, func()) {
	t.Helper()

	arquivoTemporario, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("não foi possível criar arquivo temporário %v", err)
	}

	arquivoTemporario.Write([]byte(dadosIniciais))

	removerArquivo := func() {
		arquivoTemporario.Close()
		os.Remove(arquivoTemporario.Name())
	}

	return arquivoTemporario, removerArquivo
}

func TestArmazenamentoSistemaArquivo(t *testing.T) {

	t.Run("Liga ordenada", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := poquer.NovoSistemaArquivoArmazenamentoJogador(baseDeDados)

		verificaSemErro(t, err)

		obtido := armazenamento.ObterLiga()

		esperado := []poquer.Jogador{
			{Nome: "Chris", Vitorias: 33},
			{Nome: "Cleo", Vitorias: 10},
		}

		verificaLiga(t, obtido, esperado)

		// ler de novo
		obtido = armazenamento.ObterLiga()
		verificaLiga(t, obtido, esperado)
	})

	t.Run("obter pontuação do jogador", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := poquer.NovoSistemaArquivoArmazenamentoJogador(baseDeDados)

		verificaSemErro(t, err)

		obtido := armazenamento.ObtemPontuacaoDoJogador("Chris")
		esperado := 33
		verificaPontuaçõesIguais(t, obtido, esperado)
	})

	t.Run("armazenamento de vitória para jogadores existentes", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := poquer.NovoSistemaArquivoArmazenamentoJogador(baseDeDados)

		verificaSemErro(t, err)

		armazenamento.GravarVitoria("Chris")

		obtido := armazenamento.ObtemPontuacaoDoJogador("Chris")
		esperado := 34
		verificaPontuaçõesIguais(t, obtido, esperado)
	})

	t.Run("armazenamento de vitória para jogadores existentes", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := poquer.NovoSistemaArquivoArmazenamentoJogador(baseDeDados)

		verificaSemErro(t, err)

		armazenamento.GravarVitoria("Pepper")

		obtido := armazenamento.ObtemPontuacaoDoJogador("Pepper")
		esperado := 1
		verificaPontuaçõesIguais(t, obtido, esperado)
	})

	t.Run("funciona com um arquivo vazio", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, "")
		defer limparBaseDeDados()

		_, err := poquer.NovoSistemaArquivoArmazenamentoJogador(baseDeDados)

		verificaSemErro(t, err)
	})
}

func verificaPontuaçõesIguais(t *testing.T, obtido, esperado int) {
	t.Helper()
	if obtido != esperado {
		t.Errorf("obtido %d esperado %d", obtido, esperado)
	}
}

func verificaSemErro(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("não esperava um erro mas obteve um, %v", err)
	}
}
