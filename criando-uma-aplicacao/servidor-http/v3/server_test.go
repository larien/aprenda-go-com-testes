package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type EsbocoJogadorArmazenamento struct {
	pontuacoes map[string]int
}

func (s *EsbocoJogadorArmazenamento) ObterPontuacaoJogador(nome string) int {
	pontuacao := s.pontuacoes[nome]
	return pontuacao
}

func TestObterJogadores(t *testing.T) {
	armazenamento := EsbocoJogadorArmazenamento{
		map[string]int{
			"Maria": 20,
			"Pedro": 10,
		},
	}
	servidor := &JogadorServidor{&armazenamento}

	t.Run("obter pontuação de Maria", func(t *testing.T) {
		requisicao := novaRequisicaoPontuacaoGet("Maria")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		verificarCorpoRequisicao(t, resposta.Body.String(), "20")
	})

	t.Run("obter pontuação de Pedro", func(t *testing.T) {
		requisicao := novaRequisicaoPontuacaoGet("Pedro")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		verificarCorpoRequisicao(t, resposta.Body.String(), "10")
	})

	t.Run("obter código de status 404 para jogadores não encontrados", func(t *testing.T) {
		requisicao := novaRequisicaoPontuacaoGet("Joana")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusNotFound)
	})
}

func TestArmazenamentoVitorias(t *testing.T) {
	armazenamento := EsbocoJogadorArmazenamento{
		map[string]int{},
	}
	servidor := &JogadorServidor{&armazenamento}

	t.Run("retorna status 'aceito' para chamadas ao método POST", func(t *testing.T) {
		requisicao, _ := http.NewRequest(http.MethodPost, "/jogadores/Maria", nil)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusAccepted)
	})
}

func verificarRespostaCodigoStatus(t *testing.T, recebido, esperado int) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("não recebi o codigo de status HTTP esperado, recebido %d, esperado %d", recebido, esperado)
	}
}

func novaRequisicaoPontuacaoGet(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func verificarCorpoRequisicao(t *testing.T, recebido, esperado string) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("corpo da requisição é inválido, recebido '%s' esperado '%s'", recebido, esperado)
	}
}
