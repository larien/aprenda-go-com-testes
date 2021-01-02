package poquer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGravaVitoriasEAsRetorna(t *testing.T) {
	baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[]`)
	defer limparBaseDeDados()
	armazenamento, err := NovoSistemaArquivoArmazenamentoJogador(baseDeDados)

	verificaSemErro(t, err)

	servidor := deveFazerServidorJogador(t, armazenamento)
	jogador := "Pepper"

	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))

	t.Run("obter pontuação", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoObterPontuacao(jogador))
		verificaStatus(t, resposta, http.StatusOK)

		verificaCorpoDaResposta(t, resposta.Body.String(), "3")
	})

	t.Run("obter Liga", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoDeLiga())
		verificaStatus(t, resposta, http.StatusOK)

		obtido := obterLigaDaResposta(t, resposta.Body)
		esperado := []Jogador{
			{"Pepper", 3},
		}
		verificaLiga(t, obtido, esperado)
	})
}
