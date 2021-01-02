package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGravaVitoriasEAsRetorna(t *testing.T) {
	armazenamento := NovoArmazenamentoDeJogadorNaMemoria()
	servidor := NovoServidorJogador(armazenamento)
	jogador := "Pepper"

	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))

	t.Run("obter pontuação", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoObterPontuacao(jogador))
		verificaStatus(t, resposta.Code, http.StatusOK)

		verificaCorpoDaResposta(t, resposta.Body.String(), "3")
	})

	t.Run("obter liga", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoDeLiga())
		verificaStatus(t, resposta.Code, http.StatusOK)

		obtido := obterLigaDaResposta(t, resposta.Body)
		esperado := []Jogador{
			{"Pepper", 3},
		}
		verificaLiga(t, obtido, esperado)
	})
}
