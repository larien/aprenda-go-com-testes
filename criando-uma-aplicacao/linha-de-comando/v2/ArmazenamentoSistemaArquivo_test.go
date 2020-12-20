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
		t.Fatalf("nao foi possivel criar o arquivo temporario %v", err)
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

		// Ler novamente
		obtido = armazenamento.ObterLiga()
		verificaLiga(t, obtido, esperado)
	})

	t.Run("encontrar os pontos do jogador", func(t *testing.T) {
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

	t.Run("armazenar vitórias para jogadores existentes", func(t *testing.T) {
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

	t.Run("armazenar vitórias para jogadores existentes", func(t *testing.T) {
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

	t.Run("trabalhar com um arquivo vazio", func(t *testing.T) {
		baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, "")
		defer limparBaseDeDados()

		_, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

		verificaSemErro(t, err)
	})
}

func verificaPontuacaoIgual(t *testing.T, obtido, esperado int) {
	t.Helper()
	if obtido != esperado {
		t.Errorf("recebi %d esperava %d", obtido, esperado)
	}
}

func verificaSemErro(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("não esperava um erro mas recebi um, %v", err)
	}
}
