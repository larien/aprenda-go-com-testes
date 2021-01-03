package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGravaVitoriasEAsRetorna(t *testing.T) {
	armazenamento := NovoArmazenamentoDeJogadorNaMemoria()
	servidor := ServidorJogador{armazenamento}
	jogador := "Pepper"

	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))

	resposta := httptest.NewRecorder()
	servidor.ServeHTTP(resposta, novaRequisicaoObterPontuacao(jogador))
	verificaStatus(t, resposta.Code, http.StatusOK)

	verificaCorpoDaResposta(t, resposta.Body.String(), "3")
}
