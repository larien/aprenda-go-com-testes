package poquer

import (
	"io/ioutil"
	"os"
	"testing"
)

func criarArquivoTemporario(t *testing.T, dadosIniciais string) (*os.File, func()) {
	t.Helper()

	arquivoTemporario, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("não foi possível criar arquivo temporário: %v", err)
	}

	arquivoTemporario.Write([]byte(dadosIniciais))

	removerArquivo := func() {
		os.Remove(arquivoTemporario.Name())
	}

	return arquivoTemporario, removerArquivo
}

func TestArmazenaSistemaDeArquivo(t *testing.T) {

	t.Run("liga ordenada", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "ChamadasDeVitoria": 10},
			{"Nome": "Chris", "ChamadasDeVitoria": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

		verificaSemErro(t, err)

		obtido := armazenamento.ObterLiga()

		esperado := []Jogador{
			{"Chris", 33},
			{"Cleo", 10},
		}

		verificaLiga(t, obtido, esperado)

		// ler de novo
		obtido = armazenamento.ObterLiga()
		verificaLiga(t, obtido, esperado)
	})

	t.Run("obter pontuação de jogador", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "ChamadasDeVitoria": 10},
			{"Nome": "Chris", "ChamadasDeVitoria": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

		verificaSemErro(t, err)

		obtido := armazenamento.ObterPontuacaoDeJogador("Chris")
		esperado := 33
		verificaPontuacaoIgual(t, obtido, esperado)
	})

	t.Run("vitórias armazenadas por jogadores existentes", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "ChamadasDeVitoria": 10},
			{"Nome": "Chris", "ChamadasDeVitoria": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

		verificaSemErro(t, err)

		armazenamento.GravarVitoria("Chris")

		obtido := armazenamento.ObterPontuacaoDeJogador("Chris")
		esperado := 34
		verificaPontuacaoIgual(t, obtido, esperado)
	})

	t.Run("vitórias armazenadas por jogadores existentes", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[
			{"Nome": "Cleo", "ChamadasDeVitoria": 10},
			{"Nome": "Chris", "ChamadasDeVitoria": 33}]`)
		defer limparBaseDeDados()

		armazenamento, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

		verificaSemErro(t, err)

		armazenamento.GravarVitoria("Pepper")

		obtido := armazenamento.ObterPontuacaoDeJogador("Pepper")
		esperado := 1
		verificaPontuacaoIgual(t, obtido, esperado)
	})

	t.Run("funciona com um arquivo vazio", func(t *testing.T) {
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
		t.Fatalf("não esperava um erro, mas obteve um: %v", err)
	}
}
