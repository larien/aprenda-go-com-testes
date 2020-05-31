package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistrarVitoriasEBuscarEstasVitorias(t *testing.T) {
	armazenamento := CriarJogadorArmazenamentoNaMemoria()
	servidor := JogadorServidor{armazenamento}
	jogador := "Maria"

	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoPontuacaoPost(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoPontuacaoPost(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoPontuacaoPost(jogador))

	resposta := httptest.NewRecorder()
	servidor.ServeHTTP(resposta, novaRequisicaoPontuacaoGet(jogador))
	verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)

	verificarCorpoRequisicao(t, resposta.Body.String(), "3")
}
