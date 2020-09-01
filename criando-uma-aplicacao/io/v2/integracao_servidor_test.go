package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	armazenamento := NovoArmazenamentoDeJogadorNaMemoria()
	servidor := NovoServidorDoJogador(armazenamento)
	jogador := "Pepper"

	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoPostaVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoPostaVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoPostaVitoria(jogador))

	t.Run("pega armazenamento", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, newGetScoreRequest(jogador))
		defineStatus(t, resposta.Code, http.StatusOK)

		defineCorpoDeReposta(t, resposta.Body.String(), "3")
	})

	t.Run("pega liga", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, requisicaoNovaLiga())
		defineStatus(t, resposta.Code, http.StatusOK)

		recebido := pegaLigaDaResposta(t, resposta.Body)
		esperado := []Jogador{
			{"Pepper", 3},
		}
		defineLiga(t, recebido, esperado)
	})
}
